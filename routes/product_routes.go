package routes

import (
	"github.com/denizbarcak/havuzMavisi/controllers"
	"github.com/denizbarcak/havuzMavisi/middleware"
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App) {
	// Sadece admin token ile ürün ekleyebilir, güncelleyebilir, silebilir
	app.Post("/api/products", middleware.RequireAuth(), controllers.CreateProduct)
	app.Put("/api/products/:id", middleware.RequireAuth(), controllers.UpdateProduct)
	app.Delete("/api/products/:id", middleware.RequireAuth(), controllers.DeleteProduct)
	
	// Herkes tarafından erişilebilen ürün endpoint'leri
	app.Get("/api/products", controllers.GetAllProducts) // ?category=chemicals şeklinde filtrelenebilir
	app.Get("/api/products/:id", controllers.GetProductByID)
}