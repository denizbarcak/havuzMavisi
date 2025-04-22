package main

import (
	"github.com/denizbarcak/havuzMavisi/config"
	"github.com/denizbarcak/havuzMavisi/middleware"
	"github.com/denizbarcak/havuzMavisi/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// CORS middleware'ini kur
	middleware.SetupCORS(app)
	
	config.ConnectDB()
	routes.UserRoutes(app)
	routes.CartRoutes(app)
	routes.ProductRoutes(app)
	routes.SubCategoryRoutes(app)
	routes.SetupFavoriteRoutes(app)
	routes.AdminRoutes(app)
	
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Havuz Mavisi API'ye Hoş Geldiniz!")
	})

	// 👇 Korumalı örnek route
	app.Get("/api/secret", middleware.RequireAuth(), func(c *fiber.Ctx) error {
		user := c.Locals("email")
		return c.JSON(fiber.Map{
			"message": "Giriş yaptınız",
			"user":    user,
		})
	})

	app.Listen(":8080")
}