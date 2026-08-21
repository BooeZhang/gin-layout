package role

import (
	"context"
	"errors"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gin-layout/internal/infra"
	"gin-layout/internal/page"
	"gin-layout/internal/sysuser"
)

type RoleRepository struct {
	*infra.CRUDRepository[Role, int64]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db:             db,
		CRUDRepository: infra.NewCRUDRepository[Role, int64](db),
	}
}

type roleListQuery struct {
	page.Request
	Name    *string
	Code    *string
	Enabled *bool
}

func (r *RoleRepository) List(ctx context.Context, q roleListQuery) ([]Role, int64, error) {
	var total int64
	query := q.Normalize()

	filter := r.db.Model(&Role{})
	if q.Name != nil {
		filter = filter.Where("name LIKE ?", *q.Name+"%")
	}
	if q.Code != nil {
		filter = filter.Where("code LIKE ?", *q.Code+"%")
	}
	if q.Enabled != nil {
		filter = filter.Where("enabled = ?", *q.Enabled)
	}

	if err := filter.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []Role
	err := filter.Offset(query.Offset()).Limit(query.PageSize).Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	return models, total, nil
}

func (r *RoleRepository) CreateWithMenu(ctx context.Context, role *Role, menuIDs []int64) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[Role](tx).Create(ctx, role); err != nil {
			return err
		}
		return r.replaceRoleMenu(ctx, tx, role.ID, menuIDs)
	})
	return err
}

func (r *RoleRepository) UpdateWithMenu(ctx context.Context, role *Role, menuIDs []int64) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(role).Error; err != nil {
			return err
		}
		return r.replaceRoleMenu(ctx, tx, role.ID, menuIDs)
	})
	return err
}

func (r *RoleRepository) replaceRoleMenu(ctx context.Context, tx *gorm.DB, roleID int64, menuIDs []int64) error {
	if err := tx.Where("role_id = ?", roleID).Delete(&RoleMenu{}).Error; err != nil {
		return err
	}
	if len(menuIDs) == 0 {
		return nil
	}

	items := lo.Map(menuIDs, func(item int64, _ int) RoleMenu {
		return RoleMenu{RoleID: roleID, MenuID: item}
	})
	return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&items).Error
}

func (r *RoleRepository) DeleteWithRoleID(ctx context.Context, roleID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&RoleMenu{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", roleID).Delete(&sysuser.SysUserRole{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", roleID).Delete(&Role{}).Error
	})
}

func (r *RoleRepository) FindByCode(ctx context.Context, code string) (*Role, bool, error) {
	m, err := gorm.G[Role](r.db).Where("code = ?", code).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &m, true, nil
}

func (r *RoleRepository) FindCodesByIDs(ctx context.Context, roleIDs []int64) ([]string, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}
	var codes []string
	err := r.db.WithContext(ctx).Model(&Role{}).Where("id IN ?", roleIDs).Pluck("code", &codes).Error
	if err != nil {
		return nil, err
	}
	return codes, nil
}

func (r *RoleRepository) FindByUserIDs(ctx context.Context, userIDs []int64, enabled *bool) ([]Role, error) {
	if len(userIDs) == 0 {
		return []Role{}, nil
	}

	db := r.db.WithContext(ctx).Model(&Role{}).
		Joins("JOIN sys_user_roles ON sys_user_roles.role_id = roles.id").
		Where("sys_user_roles.user_id IN ?", userIDs)
	if enabled != nil {
		db = db.Where("roles.enabled = ?", *enabled)
	}

	var models []Role
	if err := db.Order("roles.sort ASC, roles.id ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *RoleRepository) FindByIDWithPerm(ctx context.Context, roleID int64) (*Role, bool, error) {
	m, found, err := r.FindByID(ctx, roleID)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}

	var menuIDs []int64
	if err := r.db.WithContext(ctx).Model(&RoleMenu{}).Where("role_id = ?", roleID).Pluck("menu_id", &menuIDs).Error; err != nil {
		return nil, false, err
	}
	m.MenuIDs = menuIDs
	return m, true, nil
}

func (r *RoleRepository) ListAll(ctx context.Context) ([]Role, error) {
	return gorm.G[Role](r.db).Where("enabled = ?", true).Order("sort ASC, id ASC").Find(ctx)
}

func (r *RoleRepository) RoleAddUser(ctx context.Context, data []sysuser.SysUserRole) error {
	if len(data) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&data).Error
}

func (r *RoleRepository) RoleRemoveUser(ctx context.Context, roleID int64, userIDs []int64) error {
	if len(userIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Where("role_id = ? AND user_id IN ?", roleID, userIDs).Delete(&sysuser.SysUserRole{}).Error
}

func (r *RoleRepository) ReplaceRoleMenus(ctx context.Context, roleID int64, menuIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&RoleMenu{}).Error; err != nil {
			return err
		}
		if len(menuIDs) == 0 {
			return nil
		}
		items := lo.Map(menuIDs, func(item int64, _ int) RoleMenu {
			return RoleMenu{RoleID: roleID, MenuID: item}
		})
		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&items).Error
	})
}
