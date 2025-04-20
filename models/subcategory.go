package models

import "time"

type SubCategory struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
	Name       string    `json:"name" bson:"name"`
	ParentCategory string `json:"parent_category" bson:"parent_category"`
	Slug       string    `json:"slug" bson:"slug"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	CreatedBy  string    `json:"created_by" bson:"created_by"` // Admin kullanıcı ID'si
} 