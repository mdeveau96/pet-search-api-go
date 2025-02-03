package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID   `bson:"_id" json:"_id"`
	User      User                 `bson:"user" json:"user"`
	Content   string               `bson:"content" json:"content"`
	Likes     []primitive.ObjectID `bson:"likes" json:"likes"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}
