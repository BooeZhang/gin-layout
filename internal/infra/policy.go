package infra

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type roleCodeFinder interface {
	FindCodesByIDs(ctx context.Context, ids []int64) ([]string, error)
}

type CasbinManager struct {
	db        *gorm.DB
	roles     roleCodeFinder
	modelPath string
	enforcer  *casbin.TransactionalEnforcer
	logger    *zerolog.Logger
}

func NewCasbinManager(db *gorm.DB, roles roleCodeFinder, modelPath string, logger *zerolog.Logger) (*CasbinManager, error) {
	adapter, err := gormadapter.NewTransactionalAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("init casbin gorm adapter: %w", err)
	}

	enforcer, err := casbin.NewTransactionalEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("init casbin enforcer: %w", err)
	}

	c := &CasbinManager{
		db:        db,
		roles:     roles,
		modelPath: modelPath,
		enforcer:  enforcer,
		logger:    logger,
	}

	c.enforcer.LoadPolicy()
	return c, nil
}

func (m *CasbinManager) SyncUserRoles(ctx context.Context, userAccount string, roleCodes []string) error {
	err := m.enforcer.WithTransaction(ctx, func(tx *casbin.Transaction) error {
		if _, err := m.enforcer.DeleteRolesForUser(userAccount); err != nil {
			return fmt.Errorf("delete user roles: %w", err)
		}
		if _, err := m.enforcer.AddRolesForUser(userAccount, roleCodes); err != nil {
			return fmt.Errorf("add user roles: %w", err)
		}
		return nil
	})
	return err
}

func (m *CasbinManager) SyncUserRolesByIDs(ctx context.Context, userAccount string, roleIDs []int64) error {
	if len(roleIDs) == 0 {
		return m.SyncUserRoles(ctx, userAccount, nil)
	}
	if m.roles == nil {
		return fmt.Errorf("role repository is nil")
	}

	roleCodes, err := m.roles.FindCodesByIDs(ctx, roleIDs)
	if err != nil {
		return err
	}
	return m.SyncUserRoles(ctx, userAccount, roleCodes)
}

func (m *CasbinManager) SyncRolePermissions(ctx context.Context, roleCode string, permissions [][]string) error {
	err := m.enforcer.WithTransaction(ctx, func(tx *casbin.Transaction) error {
		if _, err := m.enforcer.DeletePermissionsForUser(roleCode); err != nil {
			return fmt.Errorf("delete role permissions: %w", err)
		}

		for _, permission := range permissions {
			if _, err := m.enforcer.AddPermissionsForUser(roleCode, permission); err != nil {
				return fmt.Errorf("add role permissions (role=%s, permission=%v): %w", roleCode, permission, err)
			}
		}
		return nil
	})
	return err
}

func (m *CasbinManager) AddRoleToUser(ctx context.Context, userAccount string, roleCode string) error {
	if _, err := m.enforcer.AddRoleForUser(userAccount, roleCode); err != nil {
		return fmt.Errorf("add role to user: %w", err)
	}
	return nil
}

func (m *CasbinManager) DeleteRole(ctx context.Context, roleCode string) error {
	if _, err := m.enforcer.DeleteRole(roleCode); err != nil {
		return fmt.Errorf("delete role: %w", err)
	}
	return nil
}

func (m *CasbinManager) Enforce(subject, object, action string) (bool, error) {
	return m.enforcer.Enforce(subject, object, action)
}
