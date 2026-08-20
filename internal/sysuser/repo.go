package sysuser

import (
	"context"
	"errors"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gin-layout/internal/domain"
	"gin-layout/internal/infra"
	"gin-layout/internal/page"
)

type SysUserModel struct {
	ID           int64      `gorm:"primaryKey"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
	Account      string     `gorm:"uniqueIndex;size:50;not null"`
	PasswordHash string     `gorm:"column:password;size:255;not null"`
	NickName     string     `gorm:"size:100"`
	Email        string     `gorm:"size:100;index"`
	Phone        string     `gorm:"size:20;index"`
	Avatar       string     `gorm:"size:500"`
	Enabled      bool       `gorm:"default:true"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at"`
}

func (SysUserModel) TableName() string { return "sys_user" }

func (m SysUserModel) toDomain() domain.SysUser {
	return domain.SysUser{
		ID:           m.ID,
		Account:      m.Account,
		PasswordHash: m.PasswordHash,
		NickName:     m.NickName,
		Email:        m.Email,
		Phone:        m.Phone,
		Avatar:       m.Avatar,
		Enabled:      m.Enabled,
		LastLoginAt:  m.LastLoginAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func toUserModel(u *domain.SysUser) SysUserModel {
	return SysUserModel{
		ID:           u.ID,
		Account:      u.Account,
		PasswordHash: u.PasswordHash,
		NickName:     u.NickName,
		Email:        u.Email,
		Phone:        u.Phone,
		Avatar:       u.Avatar,
		Enabled:      u.Enabled,
		LastLoginAt:  u.LastLoginAt,
	}
}

type SysUserRoleModel struct {
	UserID int64 `gorm:"primaryKey"`
	RoleID int64 `gorm:"primaryKey;index"`
}

func (SysUserRoleModel) TableName() string { return "sys_user_roles" }

type SysUserRepository struct {
	*infra.CRUDRepository[domain.SysUser, SysUserModel, int64]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *SysUserRepository {
	mapper := infra.Mapper[domain.SysUser, SysUserModel]{
		ToModel:  toUserModel,
		ToDomain: SysUserModel.toDomain,
	}
	return &SysUserRepository{
		db:             db,
		CRUDRepository: infra.NewCRUDRepository[domain.SysUser, SysUserModel, int64](db, mapper),
	}
}

type userListQuery struct {
	page.Request
	Account  *string
	NickName *string
	Email    *string
	Phone    *string
	Enabled  *bool
}

func (r *SysUserRepository) loadRoleIDs(ctx context.Context, userID int64) ([]int64, error) {
	var ids []int64
	if err := r.db.WithContext(ctx).Model(&SysUserRoleModel{}).Where("user_id = ?", userID).Pluck("role_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *SysUserRepository) List(ctx context.Context, q userListQuery) ([]domain.SysUser, int64, error) {
	var total int64
	query := q.Normalize()

	filter := r.db.Model(&SysUserModel{})
	if q.Account != nil {
		filter = filter.Where("account LIKE ?", *q.Account+"%")
	}
	if q.NickName != nil {
		filter = filter.Where("nick_name LIKE ?", *q.NickName+"%")
	}
	if q.Email != nil {
		filter = filter.Where("email LIKE ?", *q.Email+"%")
	}
	if q.Phone != nil {
		filter = filter.Where("phone LIKE ?", *q.Phone+"%")
	}
	if q.Enabled != nil {
		filter = filter.Where("enabled = ?", *q.Enabled)
	}

	if err := filter.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []SysUserModel
	err := filter.Offset(query.Offset()).Limit(query.PageSize).Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	users := lo.Map(models, func(m SysUserModel, _ int) domain.SysUser { return m.toDomain() })
	return users, total, nil
}

func (r *SysUserRepository) FindByAccount(ctx context.Context, account string) (*domain.SysUser, bool, error) {
	m, err := gorm.G[SysUserModel](r.db).Where("account = ?", account).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	u := m.toDomain()
	return &u, true, nil
}

func (r *SysUserRepository) CreateWithRoles(ctx context.Context, u *domain.SysUser, roleIDs []int64) error {
	m := toUserModel(u)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[SysUserModel](tx).Create(ctx, &m); err != nil {
			return err
		}
		u.ID = m.ID
		return r.replaceUserRoles(ctx, tx, m.ID, roleIDs)
	})
}

func (r *SysUserRepository) UpdateWithRoles(ctx context.Context, u *domain.SysUser, roleIDs []int64) error {
	m := toUserModel(u)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&m).Error; err != nil {
			return err
		}
		return r.replaceUserRoles(ctx, tx, m.ID, roleIDs)
	})
}

func (r *SysUserRepository) replaceUserRoles(ctx context.Context, tx *gorm.DB, userID int64, roleIDs []int64) error {
	if _, err := gorm.G[SysUserRoleModel](tx).Where("user_id = ?", userID).Delete(ctx); err != nil {
		return err
	}

	items := lo.Map(roleIDs, func(roleID int64, _ int) SysUserRoleModel {
		return SysUserRoleModel{UserID: userID, RoleID: roleID}
	})

	if len(items) > 0 {
		if err := tx.WithContext(ctx).Create(&items).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *SysUserRepository) FindByIDWithRoles(ctx context.Context, id int64) (*domain.SysUser, bool, error) {
	var m SysUserModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	u := m.toDomain()
	roleIDs, err := r.loadRoleIDs(ctx, id)
	if err != nil {
		return nil, false, err
	}
	u.RoleIDs = roleIDs
	return &u, true, nil
}

func (r *SysUserRepository) UpdateLastLogin(ctx context.Context, userID int64, lastLoginAt time.Time) error {
	_, err := gorm.G[SysUserModel](r.db).Where("id = ?", userID).Update(ctx, "last_login_at", lastLoginAt)
	return err
}

func (r *SysUserRepository) ReplaceUserRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	if len(roleIDs) == 0 {
		return nil
	}

	items := lo.Map(roleIDs, func(roleID int64, _ int) SysUserRoleModel {
		return SysUserRoleModel{UserID: userID, RoleID: roleID}
	})

	return r.db.Transaction(func(tx *gorm.DB) error {
		if _, err := gorm.G[SysUserRoleModel](tx).Where("user_id = ?", userID).Delete(ctx); err != nil {
			return err
		}
		return tx.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&items).Error
	})
}
