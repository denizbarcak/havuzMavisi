package routes

import (
	"github.com/denizbarcak/havuzMavisi/controllers"
	"github.com/denizbarcak/havuzMavisi/middleware"
	"github.com/gofiber/fiber/v2"
)

func AdminRoutes(app *fiber.App) {
	// Admin ürün yönetimi route'ları
	admin := app.Group("/api/admin")
	admin.Use(middleware.RequireAuth)
	admin.Use(middleware.RequireAdmin)

	// Ürün yönetimi
	admin.Post("/products", controllers.CreateProduct)
	admin.Put("/products/:id", controllers.UpdateProduct)
	admin.Delete("/products/:id", controllers.DeleteProduct)

	// İleride eklenecek diğer admin route'ları buraya eklenebilir
	// Örnek: Kullanıcı yönetimi, sipariş yönetimi, vs.
} 