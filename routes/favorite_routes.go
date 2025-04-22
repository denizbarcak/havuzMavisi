package routes

import (
	"github.com/denizbarcak/havuzMavisi/controllers"
	"github.com/denizbarcak/havuzMavisi/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupFavoriteRoutes sets up routes for favorite operations
func SetupFavoriteRoutes(app *fiber.App) {
	favorites := app.Group("/api/favorites")

	// Tüm favoriler için auth middleware gerekli
	favorites.Use(middleware.RequireAuth())

	// Favori ürün ekleme
	favorites.Post("/", controllers.AddToFavorites)

	// Favori ürün silme
	favorites.Delete("/:id", controllers.RemoveFromFavorites)

	// Favorileri listeleme (sadece ID'ler)
	favorites.Get("/", controllers.GetFavorites)

	// Favori ürünleri detaylı listeleme
	favorites.Get("/products", controllers.GetFavoriteProducts)

	// Ürünün favori olup olmadığını kontrol etme
	favorites.Get("/check/:product_id", controllers.IsFavorite)
} 