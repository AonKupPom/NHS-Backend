package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `form:"_id,omitempty" json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `form:"name" json:"name" bson:"name"`
	Type        string             `form:"type" json:"type" bson:"type"`
	Description string             `form:"description" json:"description" bson:"description"`
	Image       string             `form:"image" json:"image" bson:"image"`
	Stock       int                `form:"stock" json:"stock" bson:"stock"`
}
