package controllers

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/denizbarcak/havuzMavisi/config"
	"github.com/denizbarcak/havuzMavisi/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddToCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Giriş yapmanız gerekiyor",
		})
	}

	type AddToCartInput struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	var input AddToCartInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz veri",
		})
	}

	// Ürünün stok durumunu kontrol et
	productsCollection := config.GetCollection("products")
	var product models.Product
	err := productsCollection.FindOne(context.Background(), bson.M{"_id": input.ProductID}).Decode(&product)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Ürün bulunamadı",
		})
	}

	// Stok kontrolü
	if product.Stock < input.Quantity {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Yetersiz stok",
		})
	}

	cartItem := models.CartItem{
		ID:        uuid.New().String(),
		UserID:    userID.(string),
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
		AddedAt:   time.Now(),
	}

	collection := config.GetCollection("cart")
	_, err = collection.InsertOne(context.Background(), cartItem)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Sepete eklenemedi",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Ürün sepete eklendi",
	})
}

func GetCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Giriş yapmanız gerekiyor",
		})
	}

	collection := config.GetCollection("cart")
	filter := bson.M{"user_id": userID}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Sepet alınamadı",
		})
	}
	defer cursor.Close(context.Background())

	var cartItems []models.CartItem
	if err := cursor.All(context.Background(), &cartItems); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Veriler çözümlenemedi",
		})
	}

	return c.JSON(cartItems)
}

func GetCartItems(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	fmt.Println("Token'dan gelen user_id:", userID) // 🧪 Log ekledik

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

	cartCollection := config.GetCollection("cart")
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

	cursor, err := cartCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Sepet alınamadı"})
	}

	var cartItems []bson.M
	if err = cursor.All(ctx, &cartItems); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Sepet verileri çözümlenemedi"})
	}

	return c.Status(200).JSON(cartItems)
}

func DeleteCartItem(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Giriş yapmanız gerekiyor",
		})
	}

	itemID := c.Params("id")
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz ürün ID",
		})
	}

	collection := config.GetCollection("cart")
	filter := bson.M{
		"_id":     itemID,
		"user_id": userID,
	}

	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil || result.DeletedCount == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Sepet öğesi silinemedi",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Ürün sepetten silindi",
	})
}

func UpdateCartItem(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Giriş yapmanız gerekiyor",
		})
	}

	cartItemID := c.Params("id")
	if cartItemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz sepet öğesi ID",
		})
	}

	var updateData struct {
		Quantity int `json:"quantity"`
	}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz veri",
		})
	}

	collection := config.GetCollection("cart")
	filter := bson.M{"_id": cartItemID, "user_id": userID}
	update := bson.M{"$set": bson.M{"quantity": updateData.Quantity}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil || result.MatchedCount == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Güncelleme başarısız",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Sepet öğesi güncellendi",
	})
}

func ClearCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Giriş yapmanız gerekiyor",
		})
	}

	// Önce sepetteki ürünleri al
	cartCollection := config.GetCollection("cart")
	productsCollection := config.GetCollection("products")
	filter := bson.M{"user_id": userID}

	cursor, err := cartCollection.Find(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Sepet verileri alınamadı",
		})
	}
	defer cursor.Close(context.Background())

	var cartItems []models.CartItem
	if err := cursor.All(context.Background(), &cartItems); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Sepet verileri çözümlenemedi",
		})
	}

	// Her ürün için stok güncelleme
	for _, item := range cartItems {
		var product models.Product
		err := productsCollection.FindOne(context.Background(), bson.M{"_id": item.ProductID}).Decode(&product)
		if err != nil {
			continue // Ürün bulunamadıysa diğerine geç
		}

		// Stok güncelleme
		_, err = productsCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": item.ProductID},
			bson.M{"$set": bson.M{"stock": product.Stock - item.Quantity}},
		)
		if err != nil {
			continue // Güncelleme başarısız olduysa diğerine geç
		}
	}

	// Sepeti temizle
	_, err = cartCollection.DeleteMany(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Sepet temizlenemedi",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Ödeme başarılı ve sepet temizlendi",
	})
}

func GetCartTotal(c *fiber.Ctx) error {
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

	cartCollection := config.GetCollection("cart")
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

	groupStage := bson.D{{
		Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total", Value: bson.D{
				{Key: "$sum", Value: bson.D{
					{Key: "$multiply", Value: bson.A{"$quantity", "$product.price"}},
				}},
			}},
		},
	}}

	pipeline := mongo.Pipeline{matchStage, lookupStage, unwindStage, groupStage}

	cursor, err := cartCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Toplam fiyat hesaplanamadı"})
	}
	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil || len(result) == 0 {
		return c.JSON(fiber.Map{"total": 0})
	}

	return c.JSON(fiber.Map{
		"total": result[0]["total"],
	})
}