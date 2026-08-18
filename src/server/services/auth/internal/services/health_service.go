package services

import (
	"context"

	"gorm.io/gorm"
)

type HealthService interface {
	Check(ctx context.Context) error
}

type healthService struct {
	db *gorm.DB
}

func NewHealthService(db *gorm.DB) HealthService {
	return &healthService{
		db: db,
	}
}

func (s *healthService) Check(ctx context.Context) error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}
