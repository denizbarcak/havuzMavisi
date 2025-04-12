package routes

import (
	"github.com/denizbarcak/havuzMavisi/controllers"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	// Kullanıcı kayıt rotası
	app.Post("/api/register", controllers.RegisterUser)
	app.Post("/api/login", controllers.LoginUser)
}