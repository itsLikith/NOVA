package repository

import (
	"errors"

	"github.com/nova/auth/internal/model"
)

var ErrNotFound = errors.New("user not found")

type UserRepository interface {
	Create(user *model.User) error
	GetByUsername(username string) (*model.User, error)
	GetByID(id uint) (*model.User, error)
}
