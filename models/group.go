package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type member struct {
	UserID User   `bson:"user_id" json:"user_id"`
	Role   string `bson:"role" json:"role"`
}

type Group struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	GroupName   string             `bson:"group_name" json:"group_name"`
	Description string             `bson:"description" json:"description"`
	Members     []member           `bson:"members" json:"members"`
}
