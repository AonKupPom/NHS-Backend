package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductType struct {
	ID          primitive.ObjectID `form:"_id,omitempty" json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `form:"name" json:"name" bson:"name"`
	Description string             `form:"description" json:"description" bson:"description"`
}
