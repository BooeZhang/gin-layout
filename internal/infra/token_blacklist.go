package infra

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type TokenBlacklistModel struct {
	ID        int64     `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	TokenHash string    `gorm:"size:255;not null"`
	UserID    int64     `gorm:"index"`
	ExpiresAt time.Time `gorm:"index"`
}

func (TokenBlacklistModel) TableName() string { return "token_black_lists" }

type TokenBlacklistRepository struct {
	crud *CRUDRepository[TokenBlacklistModel, int64]
	db   *gorm.DB
}

func NewTokenBlacklistRepository(db *gorm.DB) *TokenBlacklistRepository {
	return &TokenBlacklistRepository{
		db:   db,
		crud: NewCRUDRepository[TokenBlacklistModel, int64](db),
	}
}

func (r *TokenBlacklistRepository) Exists(ctx context.Context, tokenHash string) (bool, error) {
	count, err := gorm.G[TokenBlacklistModel](r.db).
		Where("token_hash = ?", tokenHash).
		Where("expires_at > CURRENT_TIMESTAMP").
		Count(ctx, "id")
	return count > 0, NormalizeError(err)
}

func (r *TokenBlacklistRepository) Add(ctx context.Context, tokenHash string, userID int64, expiresAt time.Time) error {
	return r.crud.Create(ctx, &TokenBlacklistModel{
		TokenHash: tokenHash,
		UserID:    userID,
		ExpiresAt: expiresAt,
	})
}

// DeleteExpired removes all expired blacklist entries (expires_at < now).
func (r *TokenBlacklistRepository) DeleteExpired(ctx context.Context) error {
	_, err := gorm.G[TokenBlacklistModel](r.db).
		Where("expires_at < CURRENT_TIMESTAMP").
		Delete(ctx)
	return NormalizeError(err)
}
