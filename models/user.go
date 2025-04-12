package models

import "time"

type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"` // Hash'li saklanacak
	Role      string    `json:"role" bson:"role"`         // "user" veya "admin"
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}