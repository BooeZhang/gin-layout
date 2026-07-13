package infra

import (
	"errors"

	"gorm.io/gorm"

	"gin-layout/internal/domain"
)

func NormalizeError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrNotFound
	}
	return err
}
