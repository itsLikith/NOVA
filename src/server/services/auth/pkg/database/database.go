package database

import (
	"fmt"

	"github.com/nova/auth/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func Connect(cfg *config.Config) (*Database, error) {
	if err := ensureDatabaseExists(cfg); err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func ensureDatabaseExists(cfg *config.Config) error {
	adminDB, err := gorm.Open(postgres.Open(cfg.SystemDSN()), &gorm.Config{})
	if err != nil {
		return err
	}

	var exists bool
	if err := adminDB.Raw("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = ?)", cfg.DBName).Scan(&exists).Error; err != nil {
		return err
	}

	if !exists {
		if err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE %q", cfg.DBName)).Error; err != nil {
			return err
		}
	}

	return nil
}
