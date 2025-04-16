package controllers

import (
	"context"
	"log"
	"time"

	"github.com/denizbarcak/havuzMavisi/config"
	"github.com/denizbarcak/havuzMavisi/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz veri",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Şifre hashlenemedi",
		})
	}

	user.ID = uuid.New().String()
	user.Password = string(hashedPassword)
	user.Role = "user"
	user.CreatedAt = time.Now()

	collection := config.GetCollection("users")
	
	// Email kontrolü
	var existingUser models.User
	err = collection.FindOne(context.Background(), fiber.Map{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bu email adresi zaten kayıtlı",
		})
	}

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Println("Mongo hata:", err) // ← ekle
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kayıt başarısız",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Kayıt başarılı",
	})
}

func LoginUser(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz giriş verisi",
		})
	}

	collection := config.GetCollection("users")

	var user models.User
	err := collection.FindOne(context.Background(), fiber.Map{"email": input.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Kullanıcı bulunamadı",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Şifre hatalı",
		})
	}

	// JWT token oluştur
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 3 gün geçerli
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := []byte("havuz-mavisi-secret") // Gerçek projede .env dosyasına koy
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Token oluşturulamadı",
		})
	}

	return c.JSON(fiber.Map{
		"token": signedToken,
	})
}
