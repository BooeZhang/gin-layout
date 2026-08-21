package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// CRUDRepository 基于 GORM 的通用 CRUD,实体同时作为领域对象与持久化模型。
type CRUDRepository[M any, ID comparable] struct {
	db *gorm.DB
}

func NewCRUDRepository[M any, ID comparable](db *gorm.DB) *CRUDRepository[M, ID] {
	return &CRUDRepository[M, ID]{db: db}
}

func (r *CRUDRepository[M, ID]) Create(ctx context.Context, entity *M) error {
	if err := gorm.G[M](r.db).Create(ctx, entity); err != nil {
		return err
	}
	return nil
}

func (r *CRUDRepository[M, ID]) Update(ctx context.Context, entity *M) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *CRUDRepository[M, ID]) Delete(ctx context.Context, id ID) error {
	_, err := gorm.G[M](r.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *CRUDRepository[M, ID]) FindByID(ctx context.Context, id ID) (*M, bool, error) {
	m, err := gorm.G[M](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &m, true, nil
}

func (r *CRUDRepository[M, ID]) FindByIDs(ctx context.Context, ids []ID) ([]M, error) {
	return gorm.G[M](r.db).Where("id in (?)", ids).Find(ctx)
}
