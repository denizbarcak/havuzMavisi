package models

import "time"

type Product struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Price       float64   `json:"price" bson:"price"`
	Category    string    `json:"category" bson:"category"`
	ImageURL    string    `json:"image_url" bson:"image_url"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}