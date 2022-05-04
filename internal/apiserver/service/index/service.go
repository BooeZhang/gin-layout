package index

import (
	"context"
	"github.com/BooeZhang/gin-layout/store"
)

type IIndex interface {
	Index(ctx context.Context) error
}

type indexService struct {
	store store.Factory
	cache store.Cache
}

var _ IIndex = (*indexService)(nil)

// NewIndexService 生成index_service
func NewIndexService(store store.Factory, cache store.Cache) *indexService {
	return &indexService{store: store, cache: cache}
}
