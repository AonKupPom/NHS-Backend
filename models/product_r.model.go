package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product_rent struct {
	ID      primitive.ObjectID `form:"_id,omitempty" json:"_id,omitempty" bson:"_id,omitempty"`
	Product string             `form:"product,omitempty" json:"product,omitempty" bson:"product,omitempty"`
	Price   int                `form:"price" json:"price" bson:"price"`
}
