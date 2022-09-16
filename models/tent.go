package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tent struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Size        string             `json:"size" bson:"size"`
	Color       string             `json:"color" bson:"color"`
	Price       int                `json:"price" bson:"price"`
	Type        string             `json:"type" bson:"type"`
	Shape       string             `json:"shape" bson:"shape"`
	Description string             `json:"description" bson:"description"`
	Image       string             `json:"image" bson:"image"`
	Stock       int                `json:"stock" bson:"stock"`
}
