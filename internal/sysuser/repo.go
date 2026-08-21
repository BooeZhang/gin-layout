package sysuser

import (
	"context"
	"errors"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gin-layout/internal/infra"
	"gin-layout/internal/page"
)

type SysUserRepository struct {
	*infra.CRUDRepository[SysUser, int64]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *SysUserRepository {
	return &SysUserRepository{
		db:             db,
		CRUDRepository: infra.NewCRUDRepository[SysUser, int64](db),
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
	if err := r.db.WithContext(ctx).Model(&SysUserRole{}).Where("user_id = ?", userID).Pluck("role_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *SysUserRepository) List(ctx context.Context, q userListQuery) ([]SysUser, int64, error) {
	var total int64
	query := q.Normalize()

	filter := r.db.Model(&SysUser{})
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

	var models []SysUser
	err := filter.Offset(query.Offset()).Limit(query.PageSize).Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	return models, total, nil
}

func (r *SysUserRepository) FindByAccount(ctx context.Context, account string) (*SysUser, bool, error) {
	m, err := gorm.G[SysUser](r.db).Where("account = ?", account).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &m, true, nil
}

func (r *SysUserRepository) CreateWithRoles(ctx context.Context, u *SysUser, roleIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[SysUser](tx).Create(ctx, u); err != nil {
			return err
		}
		return r.replaceUserRoles(ctx, tx, u.ID, roleIDs)
	})
}

func (r *SysUserRepository) UpdateWithRoles(ctx context.Context, u *SysUser, roleIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(u).Error; err != nil {
			return err
		}
		return r.replaceUserRoles(ctx, tx, u.ID, roleIDs)
	})
}

func (r *SysUserRepository) replaceUserRoles(ctx context.Context, tx *gorm.DB, userID int64, roleIDs []int64) error {
	if _, err := gorm.G[SysUserRole](tx).Where("user_id = ?", userID).Delete(ctx); err != nil {
		return err
	}

	items := lo.Map(roleIDs, func(roleID int64, _ int) SysUserRole {
		return SysUserRole{UserID: userID, RoleID: roleID}
	})

	if len(items) > 0 {
		if err := tx.WithContext(ctx).Create(&items).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *SysUserRepository) FindByIDWithRoles(ctx context.Context, id int64) (*SysUser, bool, error) {
	var m SysUser
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	roleIDs, err := r.loadRoleIDs(ctx, id)
	if err != nil {
		return nil, false, err
	}
	m.RoleIDs = roleIDs
	return &m, true, nil
}

func (r *SysUserRepository) UpdateLastLogin(ctx context.Context, userID int64, lastLoginAt time.Time) error {
	_, err := gorm.G[SysUser](r.db).Where("id = ?", userID).Update(ctx, "last_login_at", lastLoginAt)
	return err
}

func (r *SysUserRepository) ReplaceUserRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	if len(roleIDs) == 0 {
		return nil
	}

	items := lo.Map(roleIDs, func(roleID int64, _ int) SysUserRole {
		return SysUserRole{UserID: userID, RoleID: roleID}
	})

	return r.db.Transaction(func(tx *gorm.DB) error {
		if _, err := gorm.G[SysUserRole](tx).Where("user_id = ?", userID).Delete(ctx); err != nil {
			return err
		}
		return tx.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&items).Error
	})
}
