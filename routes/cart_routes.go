package routes

import (
	"github.com/denizbarcak/havuzMavisi/controllers"
	"github.com/denizbarcak/havuzMavisi/middleware"
	"github.com/gofiber/fiber/v2"
)

func CartRoutes(app *fiber.App) {
	// Sepete ürün ekle
	app.Post("/api/cart/add", middleware.RequireAuth(), controllers.AddToCart) // Sepete ürün ekle	
	// Sadece giriş yapan kullanıcı sepetini görüntüleyebilir
	app.Get("/api/cart", middleware.RequireAuth(), controllers.GetCartItems) // Sepet görüntüle

	app.Delete("/api/cart/:id", middleware.RequireAuth(), controllers.DeleteCartItem) // Sepet öğesini sil

	app.Put("/api/cart/:id", middleware.RequireAuth(), controllers.UpdateCartItem) // Sepet öğesini güncelle

	app.Delete("/api/cart", middleware.RequireAuth(), controllers.ClearCart) // Sepet temizle

	app.Get("/api/cart/total", middleware.RequireAuth(), controllers.GetCartTotal)// Sepet toplam
}