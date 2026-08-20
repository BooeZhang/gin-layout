package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Mapper 领域对象 D 与持久化模型 M 之间的双向转换
type Mapper[D, M any] struct {
	ToModel  func(*D) M
	ToDomain func(M) D
}

// CRUDRepository 基于 Mapper 的通用 CRUD,入参/出参均为领域对象 D
type CRUDRepository[D, M any, ID comparable] struct {
	db     *gorm.DB
	mapper Mapper[D, M]
}

func NewCRUDRepository[D, M any, ID comparable](db *gorm.DB, mapper Mapper[D, M]) *CRUDRepository[D, M, ID] {
	return &CRUDRepository[D, M, ID]{db: db, mapper: mapper}
}

// NewModelCRUDRepository 模型即领域对象时(无需转换)的便捷构造
func NewModelCRUDRepository[M any, ID comparable](db *gorm.DB) *CRUDRepository[M, M, ID] {
	return NewCRUDRepository[M, M, ID](db, Mapper[M, M]{
		ToModel:  func(m *M) M { return *m },
		ToDomain: func(m M) M { return m },
	})
}

func (r *CRUDRepository[D, M, ID]) Create(ctx context.Context, entity *D) error {
	m := r.mapper.ToModel(entity)
	if err := gorm.G[M](r.db).Create(ctx, &m); err != nil {
		return err
	}
	*entity = r.mapper.ToDomain(m) // 回填自增 ID、时间戳等
	return nil
}

func (r *CRUDRepository[D, M, ID]) Update(ctx context.Context, entity *D) error {
	m := r.mapper.ToModel(entity)
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *CRUDRepository[D, M, ID]) Delete(ctx context.Context, id ID) error {
	_, err := gorm.G[M](r.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *CRUDRepository[D, M, ID]) FindByID(ctx context.Context, id ID) (*D, bool, error) {
	m, err := gorm.G[M](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	d := r.mapper.ToDomain(m)
	return &d, true, nil
}

func (r *CRUDRepository[D, M, ID]) FindByIDs(ctx context.Context, ids []ID) ([]D, error) {
	ms, err := gorm.G[M](r.db).Where("id in (?)", ids).Find(ctx)
	if err != nil {
		return nil, err
	}
	ds := make([]D, 0, len(ms))
	for _, m := range ms {
		ds = append(ds, r.mapper.ToDomain(m))
	}
	return ds, nil
}
