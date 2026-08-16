package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/nova/auth/internal/service"
	"github.com/nova/auth/pkg/middleware"
	"github.com/nova/auth/pkg/response"
)

type Handler struct {
	authService service.AuthService
	jwtSecret   string
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewHandler(authService service.AuthService, jwtSecret string) *Handler {
	return &Handler{authService: authService, jwtSecret: jwtSecret}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return response.OK(c, "auth service is healthy", fiber.Map{"status": "ok", "service": "auth"})
	})

	app.Post("/api/v1/auth/login", h.Login)
	app.Post("/api/v1/auth/users", middleware.AdminRequired(h.jwtSecret), h.CreateUser)
	app.Get("/api/v1/auth/validate", h.Validate)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.ErrorMessage(c, fiber.StatusBadRequest, "invalid request", "invalid request body")
	}

	auth, err := h.authService.Authenticate(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return response.ErrorMessage(c, fiber.StatusUnauthorized, "authentication failed", "invalid username or password")
		}
		return response.ErrorMessage(c, fiber.StatusInternalServerError, "authentication failed", "unable to authenticate user")
	}

	return response.OK(c, "login successful", auth)
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*service.Claims)
	if !ok || claims == nil {
		return response.ErrorMessage(c, fiber.StatusUnauthorized, "authentication required", "missing or invalid authentication token")
	}
	if claims.Role != "admin" {
		return response.ErrorMessage(c, fiber.StatusForbidden, "access denied", "admin access required")
	}

	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.ErrorMessage(c, fiber.StatusBadRequest, "invalid request", "invalid request body")
	}

	user, err := h.authService.CreateUser(req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrValidation):
			return response.ErrorMessage(c, fiber.StatusBadRequest, "invalid request", err.Error())
		case errors.Is(err, service.ErrUserExists):
			return response.ErrorMessage(c, fiber.StatusConflict, "user already exists", err.Error())
		default:
			return response.ErrorMessage(c, fiber.StatusInternalServerError, "user creation failed", "unable to create user")
		}
	}

	return response.Created(c, "user created successfully", service.UserFromModel(user))
}

func (h *Handler) Validate(c *fiber.Ctx) error {
	token, err := service.ExtractBearerToken(c.Get("Authorization"))
	if err != nil {
		message := "invalid token"
		if errors.Is(err, service.ErrMissingToken) {
			message = "missing authorization header"
		}
		return response.ErrorMessage(c, fiber.StatusUnauthorized, "authentication required", message)
	}
	claims, err := h.authService.ValidateToken(token)
	if err != nil {
		return response.ErrorMessage(c, fiber.StatusUnauthorized, "authentication required", "invalid token")
	}

	return response.OK(c, "token is valid", fiber.Map{
		"user": fiber.Map{
			"id":       claims.UserID,
			"username": claims.Username,
			"role":     claims.Role,
		},
		"expires_at": claims.ExpiresAt.Time,
	})
}
