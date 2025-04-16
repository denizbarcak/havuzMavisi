package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// User model'ini burada tanımlayalım çünkü sadece Role field'ına ihtiyacımız var
type User struct {
	Role string `json:"role"`
}

func RequireAdmin(c *fiber.Ctx) error {
	user := c.Locals("user").(User)
	if user.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Bu işlem için admin yetkisi gereklidir",
		})
	}
	return c.Next()
} 