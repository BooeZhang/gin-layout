package infra

import (
	"context"

	"gorm.io/gorm"
)

type CRUDRepository[T any, ID comparable] struct {
	db *gorm.DB
}

func NewCRUDRepository[T any, ID comparable](db *gorm.DB) *CRUDRepository[T, ID] {
	return &CRUDRepository[T, ID]{db: db}
}

func (r *CRUDRepository[T, ID]) Create(ctx context.Context, entity *T) error {
	return NormalizeError(gorm.G[T](r.db).Create(ctx, entity))
}

func (r *CRUDRepository[T, ID]) Update(ctx context.Context, entity *T) error {
	return NormalizeError(r.db.WithContext(ctx).Save(entity).Error)
}

func (r *CRUDRepository[T, ID]) Delete(ctx context.Context, id ID) error {
	_, err := gorm.G[T](r.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return NormalizeError(err)
	}
	return nil
}

func (r *CRUDRepository[T, ID]) FindByID(ctx context.Context, id ID) (*T, error) {
	entity, err := gorm.G[T](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, NormalizeError(err)
	}
	return &entity, nil
}

func (r *CRUDRepository[T, ID]) FindByIDs(ctx context.Context, ids []ID) ([]T, error) {
	roles, err := gorm.G[T](r.db).Where("id in (?)", ids).Find(ctx)
	return roles, NormalizeError(err)
}
