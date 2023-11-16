package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductSell struct {
	ID               primitive.ObjectID `form:"_id,omitempty" json:"_id,omitempty" bson:"_id,omitempty"`
	Product          string             `form:"product,omitempty" json:"product,omitempty" bson:"product,omitempty"`
	ProductAttribute string             `form:"productAttribute,omitempty" json:"productAttribute,omitempty" bson:"productAttribute,omitempty" binding:"required"`
	Price            int                `form:"price" json:"price" bson:"price"`
}
