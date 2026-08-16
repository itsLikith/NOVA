package repository

import (
	"errors"

	"github.com/nova/auth/internal/model"
	"github.com/nova/auth/pkg/database"
	"gorm.io/gorm"
)

type PostgreSQLRepository struct {
	db *database.Database
}

func NewPostgreSQLRepository(db *database.Database) *PostgreSQLRepository {
	return &PostgreSQLRepository{db: db}
}

func (r *PostgreSQLRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *PostgreSQLRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgreSQLRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}
