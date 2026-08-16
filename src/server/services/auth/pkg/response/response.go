package response

import "github.com/gofiber/fiber/v2"

type Body struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func OK(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(Body{Message: message, Data: data})
}

func Created(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Body{Message: message, Data: data})
}

func ErrorMessage(c *fiber.Ctx, status int, message, err string) error {
	return c.Status(status).JSON(Body{Message: message, Error: err})
}
