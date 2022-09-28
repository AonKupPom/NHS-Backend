package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tent struct {
	ID          primitive.ObjectID `form:"_id,omitempty" json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `form:"name" json:"name" bson:"name"`
	Size        string             `form:"size" json:"size" bson:"size"`
	Color       string             `form:"color" json:"color" bson:"color"`
	Price       int                `form:"price" json:"price" bson:"price"`
	Type        string             `form:"type" json:"type" bson:"type"`
	Shape       string             `form:"shape" json:"shape" bson:"shape"`
	Description string             `form:"description" json:"description" bson:"description"`
	Image       string             `form:"image" json:"image" bson:"image"`
	Stock       int                `form:"stock" json:"stock" bson:"stock"`
}
