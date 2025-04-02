package middleware

import (
	"strings"

	"github.com/event-system/services"
	"github.com/gofiber/fiber/v2"
)

// میدلور محافظت که چک میکنه درخواست توکن JWT معتبر داره یا نه
func Protected(authService *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// هدر احراز هویت رو میگیره
		authorization := c.Get("Authorization")

		// چک میکنه ببینه هدر احراز هویت وجود داره یا نه
		if authorization == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		// چک میکنه ببینه هدر احراز هویت پیشوند Bearer داره یا نه
		if !strings.HasPrefix(authorization, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization format")
		}

		// توکن رو استخراج میکنه
		token := strings.TrimPrefix(authorization, "Bearer ")

		// توکن رو اعتبارسنجی میکنه
		userID, err := authService.ValidateToken(token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
		}

		// آیدی کاربر رو تو locals ذخیره میکنه
		c.Locals("userID", userID)

		return c.Next()
	}
}
