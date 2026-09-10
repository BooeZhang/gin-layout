package health

import (
	"context"
	"fmt"
)

type HealthStatus struct {
	Database string `json:"database"`
	Redis    string `json:"redis"`
}

type Service struct {
	db    Pinger
	redis Pinger
}

func NewService(db, redis Pinger) *Service {
	return &Service{db: db, redis: redis}
}

func (s *Service) Check(ctx context.Context) (res HealthStatus, err error) {
	if s.db == nil {
		return res, fmt.Errorf("database is not configured")
	}
	if s.redis == nil {
		return res, fmt.Errorf("redis is not configured")
	}

	if err := s.db.Ping(ctx); err != nil {
		return res, fmt.Errorf("database ping failed: %w", err)
	}
	if err := s.redis.Ping(ctx); err != nil {
		return res, fmt.Errorf("redis ping failed: %w", err)
	}
	return HealthStatus{Database: "ok", Redis: "ok"}, nil
}
