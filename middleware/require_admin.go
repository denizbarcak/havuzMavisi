package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// RequireAdmin checks if the user has admin role
func RequireAdmin(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Bu işlem için admin yetkisi gereklidir",
		})
	}
	return c.Next()
} 