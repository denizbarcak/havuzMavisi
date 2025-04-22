package models

import "time"

// Favorite represents a user's favorited product
type Favorite struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    string    `json:"user_id" bson:"user_id"`
	ProductID string    `json:"product_id" bson:"product_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
} 