package role

import (
	"context"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gin-layout/internal/domain"
	"gin-layout/internal/infra"
)

type RoleModel struct {
	ID          int64     `gorm:"primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Name        string    `gorm:"size:50;not null"`
	Code        string    `gorm:"uniqueIndex;size:50;not null"`
	Description string    `gorm:"size:255"`
	Sort        int       `gorm:"default:0"`
	Enabled     bool      `gorm:"default:true"`
}

func (RoleModel) TableName() string { return "roles" }

func (m RoleModel) toDomain() domain.Role {
	return domain.Role{
		ID:          m.ID,
		Name:        m.Name,
		Code:        m.Code,
		Description: m.Description,
		Sort:        m.Sort,
		Enabled:     m.Enabled,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func toRoleModel(u *domain.Role) RoleModel {
	return RoleModel{
		ID:          u.ID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Name:        u.Name,
		Code:        u.Code,
		Description: u.Description,
		Sort:        u.Sort,
		Enabled:     u.Enabled,
	}
}

type RoleMenuModel struct {
	RoleID int64 `gorm:"primaryKey"`
	MenuID int64 `gorm:"primaryKey;index"`
}

func (RoleMenuModel) TableName() string { return "role_menus" }

type userRoleModel struct {
	UserID int64 `gorm:"primaryKey"`
	RoleID int64 `gorm:"primaryKey;index"`
}

func (userRoleModel) TableName() string { return "sys_user_roles" }

type RoleRepository struct {
	*infra.CRUDRepository[domain.Role, RoleModel, int64]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *RoleRepository {
	mapper := infra.Mapper[domain.Role, RoleModel]{
		ToModel:  toRoleModel,
		ToDomain: RoleModel.toDomain,
	}

	return &RoleRepository{
		db:             db,
		CRUDRepository: infra.NewCRUDRepository[domain.Role, RoleModel, int64](db, mapper),
	}
}

type roleListQuery struct {
	domain.PageRequest
	Name    *string
	Code    *string
	Enabled *bool
}

func (r *RoleRepository) List(ctx context.Context, q roleListQuery) ([]domain.Role, int64, error) {
	var total int64
	query := q.Normalize()

	filter := r.db.Model(&RoleModel{})
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

	var models []RoleModel
	err := filter.Offset(query.Offset()).Limit(query.PageSize).Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	roles := lo.Map(models, func(m RoleModel, _ int) domain.Role { return m.toDomain() })
	return roles, total, nil
}

func (r *RoleRepository) CreateWithMenu(ctx context.Context, role *domain.Role, menuIDs []int64) error {
	m := RoleModel{
		Name: role.Name, Code: role.Code, Description: role.Description,
		Sort: role.Sort, Enabled: role.Enabled,
	}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[RoleModel](tx).Create(ctx, &m); err != nil {
			return err
		}
		role.ID = m.ID
		return r.replaceRoleMenu(ctx, tx, m.ID, menuIDs)
	})
	return infra.NormalizeError(err)
}

func (r *RoleRepository) UpdateWithMenu(ctx context.Context, role *domain.Role, menuIDs []int64) error {
	m := RoleModel{
		ID: role.ID, Name: role.Name, Code: role.Code,
		Description: role.Description, Sort: role.Sort, Enabled: role.Enabled,
	}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&m).Error; err != nil {
			return err
		}
		return r.replaceRoleMenu(ctx, tx, m.ID, menuIDs)
	})
	return infra.NormalizeError(err)
}

func (r *RoleRepository) replaceRoleMenu(ctx context.Context, tx *gorm.DB, roleID int64, menuIDs []int64) error {
	if err := tx.Where("role_id = ?", roleID).Delete(&RoleMenuModel{}).Error; err != nil {
		return err
	}
	if len(menuIDs) == 0 {
		return nil
	}

	items := lo.Map(menuIDs, func(item int64, _ int) RoleMenuModel {
		return RoleMenuModel{RoleID: roleID, MenuID: item}
	})
	return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&items).Error
}

func (r *RoleRepository) DeleteWithRoleID(ctx context.Context, roleID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&RoleMenuModel{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", roleID).Delete(&userRoleModel{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", roleID).Delete(&RoleModel{}).Error
	})
}

func (r *RoleRepository) FindByCode(ctx context.Context, code string) (*domain.Role, error) {
	m, err := gorm.G[RoleModel](r.db).Where("code = ?", code).First(ctx)
	if err != nil {
		return nil, infra.NormalizeError(err)
	}
	role := m.toDomain()
	return &role, nil
}

func (r *RoleRepository) FindCodesByIDs(ctx context.Context, roleIDs []int64) ([]string, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}
	var codes []string
	err := r.db.WithContext(ctx).Model(&RoleModel{}).Where("id IN ?", roleIDs).Pluck("code", &codes).Error
	if err != nil {
		return nil, infra.NormalizeError(err)
	}
	return codes, nil
}

func (r *RoleRepository) FindByUserIDs(ctx context.Context, userIDs []int64, enabled *bool) ([]domain.Role, error) {
	if len(userIDs) == 0 {
		return []domain.Role{}, nil
	}

	db := r.db.WithContext(ctx).Model(&RoleModel{}).
		Joins("JOIN sys_user_roles ON sys_user_roles.role_id = roles.id").
		Where("sys_user_roles.user_id IN ?", userIDs)
	if enabled != nil {
		db = db.Where("roles.enabled = ?", *enabled)
	}

	var models []RoleModel
	if err := db.Order("roles.sort ASC, roles.id ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	roles := lo.Map(models, func(m RoleModel, _ int) domain.Role { return m.toDomain() })
	return roles, nil
}

func (r *RoleRepository) FindByIDWithPerm(ctx context.Context, roleID int64) (*domain.Role, error) {
	m, err := r.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	var menuIDs []int64
	r.db.WithContext(ctx).Model(&RoleMenuModel{}).Where("role_id = ?", roleID).Pluck("menu_id", &menuIDs)
	m.MenuIDs = menuIDs
	return m, nil
}

func (r *RoleRepository) ListAll(ctx context.Context) ([]domain.Role, error) {
	models, err := gorm.G[RoleModel](r.db).Where("enabled = ?", true).Order("sort ASC, id ASC").Find(ctx)
	if err != nil {
		return nil, err
	}
	roles := lo.Map(models, func(m RoleModel, _ int) domain.Role { return m.toDomain() })
	return roles, nil
}

func (r *RoleRepository) RoleAddUser(ctx context.Context, data []domain.UserRole) error {
	if len(data) == 0 {
		return nil
	}
	items := lo.Map(data, func(item domain.UserRole, _ int) userRoleModel {
		return userRoleModel{UserID: item.UserID, RoleID: item.RoleID}
	})
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&items).Error
}

func (r *RoleRepository) RoleRemoveUser(ctx context.Context, roleID int64, userIDs []int64) error {
	if len(userIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Where("role_id = ? AND user_id IN ?", roleID, userIDs).Delete(&userRoleModel{}).Error
}

func (r *RoleRepository) ReplaceRoleMenus(ctx context.Context, roleID int64, menuIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&RoleMenuModel{}).Error; err != nil {
			return err
		}
		if len(menuIDs) == 0 {
			return nil
		}
		items := lo.Map(menuIDs, func(item int64, _ int) RoleMenuModel {
			return RoleMenuModel{RoleID: roleID, MenuID: item}
		})
		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&items).Error
	})
}
