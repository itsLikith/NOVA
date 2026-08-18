package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nova/auth/config"
	"github.com/nova/auth/internal/models"
	"github.com/nova/auth/internal/repository"
	"github.com/nova/auth/pkg/utils"

	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("userid already exists")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidUserData    = errors.New("invalid user data")
)

type LoginResult struct {
	User  *models.User
	Token string
}

type AuthService interface {
	Login(
		ctx context.Context,
		userID string,
		password string,
	) (*LoginResult, error)

	CreateUser(
		ctx context.Context,
		userID string,
		email string,
		password string,
	) (*models.User, error)

	EnsureAdmin(ctx context.Context) error
}

type authService struct {
	repository repository.AuthRepository
	config     *config.Config
}

func NewAuthService(
	repo repository.AuthRepository,
	cfg *config.Config,
) AuthService {
	return &authService{
		repository: repo,
		config:     cfg,
	}
}

func (s *authService) Login(
	ctx context.Context,
	userID string,
	password string,
) (*LoginResult, error) {
	userID = strings.TrimSpace(userID)

	if userID == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.repository.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, fmt.Errorf("find user: %w", err)
	}

	if err := utils.ComparePassword(user.PasswordHash, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(
		user.UserID,
		user.Role,
		s.config.JWTIssuer,
		s.config.JWTSecret,
		s.config.JWTExpirationHour,
	)

	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &LoginResult{
		User:  user,
		Token: token,
	}, nil
}

func (s *authService) CreateUser(
	ctx context.Context,
	userID string,
	email string,
	password string,
) (*models.User, error) {
	userID = strings.TrimSpace(userID)
	email = strings.TrimSpace(strings.ToLower(email))

	if userID == "" || email == "" || password == "" {
		return nil, ErrInvalidUserData
	}

	if len(password) < 8 {
		return nil, ErrInvalidUserData
	}

	exists, err := s.repository.ExistsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("check userid: %w", err)
	}

	if exists {
		return nil, ErrUserAlreadyExists
	}

	exists, err = s.repository.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("check email: %w", err)
	}

	if exists {
		return nil, ErrEmailAlreadyExists
	}

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		UserID:       userID,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         models.RoleUser,
	}

	if err := s.repository.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *authService) EnsureAdmin(ctx context.Context) error {
	user, err := s.repository.FindByUserID(
		ctx,
		s.config.AdminUserID,
	)

	if err == nil {
		if user.Role != models.RoleAdmin {
			return fmt.Errorf(
				"configured admin userid already exists with role %q",
				user.Role,
			)
		}

		return nil
	}

	if !errors.Is(err, repository.ErrUserNotFound) {
		return fmt.Errorf("check admin: %w", err)
	}

	passwordHash, err := utils.HashPassword(
		s.config.AdminPassword,
	)

	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	admin := &models.User{
		UserID:       s.config.AdminUserID,
		Email:        strings.ToLower(s.config.AdminEmail),
		PasswordHash: passwordHash,
		Role:         models.RoleAdmin,
	}

	if err := s.repository.Create(ctx, admin); err != nil {
		return fmt.Errorf("create admin: %w", err)
	}

	return nil
}

// Keep gorm imported available if you later want repository-level
// duplicate error translation.
var _ = gorm.ErrRecordNotFound
