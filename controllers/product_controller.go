package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/denizbarcak/havuzMavisi/config"
	"github.com/denizbarcak/havuzMavisi/models"
)

func CreateProduct(c *fiber.Ctx) error {
	// Token'dan gelen role bilgisini al
	role := c.Locals("role")
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Bu işlem sadece admin tarafından yapılabilir",
		})
	}

	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz ürün verisi",
		})
	}

	product.ID = uuid.New().String()
	product.CreatedAt = time.Now()

	collection := config.GetCollection("products")
	_, err := collection.InsertOne(context.Background(), product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ürün eklenemedi",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Ürün başarıyla eklendi",
		"product": product,
	})
}
func UpdateProduct(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Bu işlem sadece admin tarafından yapılabilir",
		})
	}

	productID := c.Params("id")
	var updatedData models.Product

	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Veri çözümlenemedi",
		})
	}

	collection := config.GetCollection("products")
	filter := bson.M{"_id": productID}
	update := bson.M{
		"$set": bson.M{
			"name":        updatedData.Name,
			"description": updatedData.Description,
			"price":       updatedData.Price,
			"category":    updatedData.Category,
			"image_url":   updatedData.ImageURL,
			"stock":       updatedData.Stock,
		},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ürün güncellenemedi",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Ürün başarıyla güncellendi",
	})
}
func DeleteProduct(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Bu işlem sadece admin tarafından yapılabilir",
		})
	}

	productID := c.Params("id")
	collection := config.GetCollection("products")

	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": productID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ürün silinemedi",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Ürün başarıyla silindi",
	})
}
func GetAllProducts(c *fiber.Ctx) error {
	collection := config.GetCollection("products")

	// URL'den kategori ve alt kategori parametrelerini al
	category := c.Query("category")
	subcategory := c.Query("subcategory")

	// Sorgu filtresini oluştur
	filter := bson.M{}
	if category != "" {
		filter["category"] = category
	}
	if subcategory != "" {
		filter["subcategory"] = subcategory
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ürünler alınamadı",
		})
	}
	defer cursor.Close(context.Background())

	var products []models.Product
	if err := cursor.All(context.Background(), &products); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ürünler çözümlenemedi",
		})
	}

	return c.JSON(products)
}

// Tek bir ürünü ID'ye göre getir
func GetProductByID(c *fiber.Ctx) error {
	productID := c.Params("id")
	collection := config.GetCollection("products")

	var product models.Product
	err := collection.FindOne(context.Background(), bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Ürün bulunamadı",
		})
	}

	return c.JSON(product)
}