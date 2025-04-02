package config

import (
	"github.com/gofiber/fiber/v2"
)

// یه هندلر خطای سفارشی برای Fiber
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	c.Status(code)
	return c.JSON(fiber.Map{
		"error":   true,
		"message": err.Error(),
	})
}
