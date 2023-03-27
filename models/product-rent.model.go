package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductRent struct {
	ID      primitive.ObjectID `form:"_id,omitempty" json:"_id,omitempty" bson:"_id,omitempty"`
	Product string             `form:"product,omitempty" json:"product,omitempty" bson:"product,omitempty" binding:"required"`
	Price   int                `form:"price" json:"price" bson:"price" binding:"required"`
}
