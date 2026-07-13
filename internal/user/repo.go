package user

import (
	"context"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gin-layout/internal/domain"
	"gin-layout/internal/infra"
)

type UserModel struct {
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

func (UserModel) TableName() string { return "users" }

func (m UserModel) toDomain() domain.User {
	return domain.User{
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

func toUserModel(u domain.User) UserModel {
	return UserModel{
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

type UserRoleModel struct {
	UserID int64 `gorm:"primaryKey"`
	RoleID int64 `gorm:"primaryKey;index"`
}

func (UserRoleModel) TableName() string { return "user_roles" }

type PGRepository struct {
	crud *infra.CRUDRepository[UserModel, int64]
	db   *gorm.DB
}

func NewRepository(db *gorm.DB) *PGRepository {
	return &PGRepository{
		db:   db,
		crud: infra.NewCRUDRepository[UserModel, int64](db),
	}
}

type userListQuery struct {
	domain.PageRequest
	Account  *string
	NickName *string
	Email    *string
	Phone    *string
	Enabled  *bool
}

func (r *PGRepository) Create(ctx context.Context, u *domain.User) error {
	m := toUserModel(*u)
	if err := r.crud.Create(ctx, &m); err != nil {
		return err
	}
	u.ID = m.ID
	u.CreatedAt = m.CreatedAt
	u.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *PGRepository) Update(ctx context.Context, u *domain.User) error {
	m := toUserModel(*u)
	return r.crud.Update(ctx, &m)
}

func (r *PGRepository) Delete(ctx context.Context, id int64) error {
	return r.crud.Delete(ctx, id)
}

func (r *PGRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	m, err := r.crud.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	u := m.toDomain()
	return &u, nil
}

func (r *PGRepository) loadRoleIDs(ctx context.Context, userID int64) ([]int64, error) {
	var ids []int64
	if err := r.db.WithContext(ctx).Model(&UserRoleModel{}).Where("user_id = ?", userID).Pluck("role_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *PGRepository) List(ctx context.Context, q userListQuery) ([]domain.User, int64, error) {
	var total int64
	query := q.Normalize()

	filter := r.db.Model(&UserModel{})
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

	var models []UserModel
	err := filter.Offset(query.Offset()).Limit(query.PageSize).Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	users := lo.Map(models, func(m UserModel, _ int) domain.User { return m.toDomain() })
	return users, total, nil
}

func (r *PGRepository) FindByAccount(ctx context.Context, account string) (*domain.User, error) {
	m, err := gorm.G[UserModel](r.db).Where("account = ?", account).First(ctx)
	if err != nil {
		return nil, infra.NormalizeError(err)
	}
	u := m.toDomain()
	return &u, nil
}

func (r *PGRepository) CreateWithRoles(ctx context.Context, u *domain.User, roleIDs []int64) error {
	m := toUserModel(*u)
	return infra.NormalizeError(r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[UserModel](tx).Create(ctx, &m); err != nil {
			return err
		}
		u.ID = m.ID
		return r.replaceUserRoles(ctx, tx, m.ID, roleIDs)
	}))
}

func (r *PGRepository) UpdateWithRoles(ctx context.Context, u *domain.User, roleIDs []int64) error {
	m := toUserModel(*u)
	return infra.NormalizeError(r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&m).Error; err != nil {
			return err
		}
		return r.replaceUserRoles(ctx, tx, m.ID, roleIDs)
	}))
}

func (r *PGRepository) replaceUserRoles(ctx context.Context, tx *gorm.DB, userID int64, roleIDs []int64) error {
	if _, err := gorm.G[UserRoleModel](tx).Where("user_id = ?", userID).Delete(ctx); err != nil {
		return err
	}

	items := lo.Map(roleIDs, func(roleID int64, _ int) UserRoleModel {
		return UserRoleModel{UserID: userID, RoleID: roleID}
	})

	if len(items) > 0 {
		if err := tx.WithContext(ctx).Create(&items).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *PGRepository) FindByIDWithRoles(ctx context.Context, id int64) (*domain.User, error) {
	var m UserModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, infra.NormalizeError(err)
	}
	u := m.toDomain()
	roleIDs, err := r.loadRoleIDs(ctx, id)
	if err != nil {
		return nil, err
	}
	u.RoleIDs = roleIDs
	return &u, nil
}

func (r *PGRepository) UpdateLastLogin(ctx context.Context, userID int64, lastLoginAt time.Time) error {
	_, err := gorm.G[UserModel](r.db).Where("id = ?", userID).Update(ctx, "last_login_at", lastLoginAt)
	return infra.NormalizeError(err)
}

func (r *PGRepository) ReplaceUserRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	if len(roleIDs) == 0 {
		return nil
	}

	items := lo.Map(roleIDs, func(roleID int64, _ int) UserRoleModel {
		return UserRoleModel{UserID: userID, RoleID: roleID}
	})

	return r.db.Transaction(func(tx *gorm.DB) error {
		if _, err := gorm.G[UserRoleModel](tx).Where("user_id = ?", userID).Delete(ctx); err != nil {
			return err
		}
		return tx.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&items).Error
	})
}

func (r *PGRepository) FindByIDs(ctx context.Context, ids []int64) ([]domain.User, error) {
	models, err := r.crud.FindByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	users := lo.Map(models, func(m UserModel, _ int) domain.User { return m.toDomain() })
	return users, nil
}
