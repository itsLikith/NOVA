package service

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/nova/auth/internal/model"
	"github.com/nova/auth/internal/repository"
)

const (
	defaultTokenTTL = 12 * time.Hour
)

var (
	ErrAdminCredentialsNotConfigured = errors.New("admin credentials are not configured")
	ErrInvalidCredentials            = errors.New("invalid username or password")
	ErrInvalidToken                  = errors.New("invalid token")
	ErrMissingToken                  = errors.New("missing token")
	ErrUserExists                    = errors.New("username already exists")
	ErrValidation                    = errors.New("username and password are required")
)

type Claims struct {
	UserID   uint       `json:"user_id"`
	Username string     `json:"username"`
	Role     model.Role `json:"role"`
	jwt.RegisteredClaims
}

type Service struct {
	repo      repository.UserRepository
	secret    string
	jwtTTL    time.Duration
	adminUser string
	adminPass string
}

type AuthService interface {
	SeedAdmin() error
	Authenticate(username, password string) (*TokenResponse, error)
	CreateUser(username, password string) (*model.User, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type TokenResponse struct {
	Token     string    `json:"token"`
	TokenType string    `json:"token_type"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserDTO   `json:"user"`
}

type UserDTO struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Role      model.Role `json:"role"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func NewService(repo repository.UserRepository, secret, adminUsername, adminPassword string) *Service {
	if secret == "" {
		secret = "nova-dev-secret"
	}

	return &Service{
		repo:      repo,
		secret:    secret,
		jwtTTL:    defaultTokenTTL,
		adminUser: adminUsername,
		adminPass: adminPassword,
	}
}

func (s *Service) SeedAdmin() error {
	if strings.TrimSpace(s.adminUser) == "" || strings.TrimSpace(s.adminPass) == "" {
		return ErrAdminCredentialsNotConfigured
	}

	existing, err := s.repo.GetByUsername(s.adminUser)
	if err == nil && existing != nil {
		return nil
	}
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(s.adminPass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &model.User{
		Username: s.adminUser,
		Password: string(hash),
		Role:     model.RoleAdmin,
	}

	return s.repo.Create(admin)
}

func (s *Service) Authenticate(username, password string) (*TokenResponse, error) {
	username = strings.TrimSpace(username)
	if username == "" || strings.TrimSpace(password) == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.repo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresAt: expiresAt,
		User:      UserFromModel(user),
	}, nil
}

func (s *Service) CreateUser(username, password string) (*model.User, error) {
	username = strings.TrimSpace(username)
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return nil, ErrValidation
	}

	existing, err := s.repo.GetByUsername(username)
	if err == nil && existing != nil {
		return nil, ErrUserExists
	}
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: string(hash),
		Role:     model.RoleUser,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	return ValidateJWT(s.secret, tokenString)
}

func ExtractBearerToken(authHeader string) (string, error) {
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return "", ErrMissingToken
	}

	parts := strings.Fields(authHeader)
	if len(parts) == 1 {
		return parts[0], nil
	}
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1], nil
	}

	return "", ErrInvalidToken
}

func ValidateJWT(secret, tokenString string) (*Claims, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, ErrMissingToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func UserFromModel(user *model.User) UserDTO {
	dto := UserDTO{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
	if !user.CreatedAt.IsZero() {
		dto.CreatedAt = &user.CreatedAt
	}
	return dto
}

func (s *Service) generateToken(user *model.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.jwtTTL).UTC()
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	return tokenString, expiresAt, err
}
