package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type URL struct {
	OwnerIP     string             `bson:"owner_ip"`
	OriginalURL string             `bson:"original_url"`
	ShortURL    string             `bson:"short_url"`
	CreatedAt   primitive.DateTime `bson:"created_at"`
}
