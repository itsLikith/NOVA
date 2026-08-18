package repository

import (
	"context"
	"errors"

	"github.com/nova/auth/internal/models"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type AuthRepository interface {
	FindByUserID(ctx context.Context, userID string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)

	Create(ctx context.Context, user *models.User) error

	ExistsByUserID(ctx context.Context, userID string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) FindByUserID(
	ctx context.Context,
	userID string,
) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepository) Create(
	ctx context.Context,
	user *models.User,
) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *authRepository) ExistsByUserID(
	ctx context.Context,
	userID string,
) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Count(&count).Error

	return count > 0, err
}

func (r *authRepository) ExistsByEmail(
	ctx context.Context,
	email string,
) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error

	return count > 0, err
}
