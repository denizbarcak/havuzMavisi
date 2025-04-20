package routes

import (
	"github.com/denizbarcak/havuzMavisi/controllers"
	"github.com/denizbarcak/havuzMavisi/middleware"
	"github.com/gofiber/fiber/v2"
)

func SubCategoryRoutes(app *fiber.App) {
	// Herkes tarafından erişilebilen alt kategori endpoint'leri
	app.Get("/api/categories/:category/subcategories", controllers.GetSubCategoriesByParent)
	app.Get("/api/subcategories", controllers.GetAllSubCategories)
	
	// Sadece admin tarafından erişilebilen alt kategori endpoint'leri
	app.Post("/api/subcategories", middleware.RequireAuth(), controllers.CreateSubCategory)
	app.Delete("/api/subcategories/:id", middleware.RequireAuth(), controllers.DeleteSubCategory)
} 