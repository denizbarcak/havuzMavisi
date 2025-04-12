package models

import "time"

type CartItem struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    string    `json:"user_id" bson:"user_id"`
	ProductID string    `json:"product_id" bson:"product_id"`
	Quantity  int       `json:"quantity" bson:"quantity"`
	AddedAt   time.Time `json:"added_at" bson:"added_at"`
}