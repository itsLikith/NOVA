package handlers

import (
	"errors"
	"time"

	"github.com/nova/auth/config"
	"github.com/nova/auth/internal/models"
	"github.com/nova/auth/internal/services"
	"github.com/nova/auth/pkg/middleware"
	"github.com/nova/pkg/logger"
	"github.com/nova/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	service services.AuthService
	cfg     *config.Config
}

func NewAuthHandler(service services.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		service: service,
		cfg:     cfg,
	}
}

type LoginRequest struct {
	UserID   string `json:"userid"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	UserID   string `json:"userid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	UserID string `json:"userid"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req LoginRequest

	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid request body",
		)
	}

	result, err := h.service.Login(
		c.Context(),
		req.UserID,
		req.Password,
	)

	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return fiber.NewError(
				fiber.StatusUnauthorized,
				"invalid credentials",
			)
		}

		return fiber.NewError(
			fiber.StatusInternalServerError,
			"internal server error",
		)
	}

	// set http-only cookie with the token
	cookie := &fiber.Cookie{
		Name:     "token",
		Value:    result.Token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   h.cfg.AppEnv != "development",
		SameSite: "Lax",
		Expires:  time.Now().Add(time.Duration(h.cfg.JWTExpirationHour) * time.Hour),
	}

	// Log cookie set event (do not log token value)
	logger.Info("Setting auth cookie for user " + result.User.UserID + " expires " + cookie.Expires.String())

	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(response.SendSuccessResponse(
		fiber.StatusOK,
		"login successful",
		fiber.Map{
			"token": result.Token,
			"user": UserResponse{
				UserID: result.User.UserID,
				Email:  result.User.Email,
				Role:   result.User.Role,
			},
		},
	))
}

func (h *AuthHandler) Validate(c fiber.Ctx) error {
	// JWTAuth middleware will populate locals
	userID, _ := c.Locals(middleware.UserIDKey).(string)
	role, _ := c.Locals(middleware.RoleKey).(string)

	if userID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	return c.Status(fiber.StatusOK).JSON(response.SendSuccessResponse(
		fiber.StatusOK,
		"token valid",
		fiber.Map{"user": UserResponse{UserID: userID, Role: role}},
	))
}

func (h *AuthHandler) CreateUser(c fiber.Ctx) error {
	var req CreateUserRequest

	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid request body",
		)
	}

	user, err := h.service.CreateUser(
		c.Context(),
		req.UserID,
		req.Email,
		req.Password,
	)

	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidUserData):
			return fiber.NewError(
				fiber.StatusBadRequest,
				"invalid user data",
			)

		case errors.Is(err, services.ErrUserAlreadyExists):
			return fiber.NewError(
				fiber.StatusConflict,
				"userid already exists",
			)

		case errors.Is(err, services.ErrEmailAlreadyExists):
			return fiber.NewError(
				fiber.StatusConflict,
				"email already exists",
			)

		default:
			return fiber.NewError(
				fiber.StatusInternalServerError,
				"internal server error",
			)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response.SendSuccessResponse(
		fiber.StatusCreated,
		"user created successfully",
		fiber.Map{"user": UserResponse{
			UserID: user.UserID,
			Email:  user.Email,
			Role:   user.Role,
		},
		},
	))
}

var _ = models.RoleUser
