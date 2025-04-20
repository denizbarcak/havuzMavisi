package controllers

import (
	"context"
	"strings"
	"time"

	"github.com/denizbarcak/havuzMavisi/config"
	"github.com/denizbarcak/havuzMavisi/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// Alt kategori oluşturma
func CreateSubCategory(c *fiber.Ctx) error {
	// Admin kontrolü
	role := c.Locals("role")
	userId := c.Locals("user_id").(string)
	
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Bu işlem sadece admin tarafından yapılabilir",
		})
	}

	var subCategory models.SubCategory

	if err := c.BodyParser(&subCategory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz alt kategori verisi",
		})
	}

	// Validasyon
	if subCategory.Name == "" || subCategory.ParentCategory == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Alt kategori adı ve üst kategori zorunludur",
		})
	}

	// Unique ID ve slug oluştur
	subCategory.ID = uuid.New().String()
	subCategory.CreatedAt = time.Now()
	subCategory.CreatedBy = userId
	
	// Slug oluşturma - boşlukları tire ile değiştir ve küçük harfe dönüştür
	subCategory.Slug = strings.ToLower(strings.ReplaceAll(subCategory.Name, " ", "-"))

	collection := config.GetCollection("subcategories")
	
	// Aynı isimde alt kategori var mı kontrolü
	var existingSubCategory models.SubCategory
	err := collection.FindOne(
		context.Background(), 
		bson.M{
			"name": subCategory.Name, 
			"parent_category": subCategory.ParentCategory,
		},
	).Decode(&existingSubCategory)
	
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bu isimde bir alt kategori zaten mevcut",
		})
	}

	_, err = collection.InsertOne(context.Background(), subCategory)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Alt kategori eklenemedi",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Alt kategori başarıyla eklendi",
		"subcategory": subCategory,
	})
}

// Bir kategoriye ait tüm alt kategorileri getir
func GetSubCategoriesByParent(c *fiber.Ctx) error {
	parentCategory := c.Params("category")
	
	if parentCategory == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Üst kategori belirtilmelidir",
		})
	}

	collection := config.GetCollection("subcategories")
	
	cursor, err := collection.Find(
		context.Background(), 
		bson.M{"parent_category": parentCategory},
	)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Alt kategoriler alınamadı",
		})
	}
	defer cursor.Close(context.Background())

	var subCategories []models.SubCategory
	if err := cursor.All(context.Background(), &subCategories); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Alt kategoriler çözümlenemedi",
		})
	}

	return c.JSON(subCategories)
}

// Tüm alt kategorileri getir
func GetAllSubCategories(c *fiber.Ctx) error {
	collection := config.GetCollection("subcategories")
	
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Alt kategoriler alınamadı",
		})
	}
	defer cursor.Close(context.Background())

	var subCategories []models.SubCategory
	if err := cursor.All(context.Background(), &subCategories); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Alt kategoriler çözümlenemedi",
		})
	}

	return c.JSON(subCategories)
}

// Alt kategori silme (sadece admin)
func DeleteSubCategory(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Bu işlem sadece admin tarafından yapılabilir",
		})
	}

	subCategoryID := c.Params("id")
	collection := config.GetCollection("subcategories")

	// Önce alt kategoriye ait ürünleri kontrol et
	productsCollection := config.GetCollection("products")
	count, err := productsCollection.CountDocuments(
		context.Background(),
		bson.M{"subcategory": subCategoryID},
	)
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Alt kategori kontrol edilemedi",
		})
	}
	
	if count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bu alt kategoriye ait ürünler mevcut. Önce ürünleri başka bir kategoriye taşıyın.",
		})
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": subCategoryID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Alt kategori silinemedi",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Alt kategori başarıyla silindi",
	})
} 