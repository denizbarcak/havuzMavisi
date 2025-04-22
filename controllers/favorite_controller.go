package controllers

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/denizbarcak/havuzMavisi/config"
	"github.com/denizbarcak/havuzMavisi/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AddToFavorites adds a product to user's favorites
func AddToFavorites(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Giriş yapmanız gerekiyor",
		})
	}

	var input struct {
		ProductID string `json:"product_id"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz veri",
		})
	}

	// Ürünün varlığını kontrol et
	productsCollection := config.GetCollection("products")
	var product models.Product
	err := productsCollection.FindOne(context.Background(), bson.M{"_id": input.ProductID}).Decode(&product)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Ürün bulunamadı",
		})
	}

	// Ürün zaten favorilerde mi kontrol et
	favoritesCollection := config.GetCollection("favorites")
	filter := bson.M{
		"user_id":    userID.(string),
		"product_id": input.ProductID,
	}

	var existingFavorite models.Favorite
	err = favoritesCollection.FindOne(context.Background(), filter).Decode(&existingFavorite)
	if err == nil {
		// Ürün zaten favorilerde
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Bu ürün zaten favorilerinizde",
		})
	}

	// Yeni favori oluştur
	favorite := models.Favorite{
		ID:        uuid.New().String(),
		UserID:    userID.(string),
		ProductID: input.ProductID,
		CreatedAt: time.Now(),
	}

	_, err = favoritesCollection.InsertOne(context.Background(), favorite)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Favorilere eklenemedi",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Ürün favorilere eklendi",
		"favorite": favorite,
	})
}

// RemoveFromFavorites removes a product from user's favorites
func RemoveFromFavorites(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Giriş yapmanız gerekiyor",
		})
	}

	favoriteID := c.Params("id")
	if favoriteID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz favori ID",
		})
	}

	collection := config.GetCollection("favorites")
	filter := bson.M{
		"_id":     favoriteID,
		"user_id": userID,
	}

	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Favori silinemedi",
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Favori bulunamadı",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Ürün favorilerden kaldırıldı",
	})
}

// GetFavorites returns all favorites for a user
func GetFavorites(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Giriş yapmanız gerekiyor",
		})
	}

	collection := config.GetCollection("favorites")
	filter := bson.M{"user_id": userID}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Favoriler alınamadı",
		})
	}
	defer cursor.Close(context.Background())

	var favorites []models.Favorite
	if err := cursor.All(context.Background(), &favorites); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Veriler çözümlenemedi",
		})
	}

	return c.JSON(favorites)
}

// GetFavoriteProducts returns all favorite products with complete product details
func GetFavoriteProducts(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Giriş yapmanız gerekiyor",
		})
	}
	userIDStr, ok := userID.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user_id tipi geçersiz",
		})
	}

	favoritesCollection := config.GetCollection("favorites")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userIDStr}}}}

	lookupStage := bson.D{{
		Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "products"},
			{Key: "localField", Value: "product_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "product"},
		},
	}}

	unwindStage := bson.D{{Key: "$unwind", Value: "$product"}}

	pipeline := mongo.Pipeline{matchStage, lookupStage, unwindStage}

	cursor, err := favoritesCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Favoriler alınamadı"})
	}

	var favoriteItems []bson.M
	if err = cursor.All(ctx, &favoriteItems); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Favori verileri çözümlenemedi"})
	}

	return c.Status(200).JSON(favoriteItems)
}

// IsFavorite checks if a product is in user's favorites
func IsFavorite(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Giriş yapmanız gerekiyor",
		})
	}

	productID := c.Params("product_id")
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz ürün ID",
		})
	}

	collection := config.GetCollection("favorites")
	filter := bson.M{
		"user_id":    userID.(string),
		"product_id": productID,
	}

	var favorite models.Favorite
	err := collection.FindOne(context.Background(), filter).Decode(&favorite)
	
	isFavorite := err == nil
	
	return c.JSON(fiber.Map{
		"is_favorite": isFavorite,
		"favorite_id": favorite.ID,
	})
} 