package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductAttribute struct {
	ID      primitive.ObjectID `form:"_id,omitempty" json:"_id,omitempty" bson:"_id,omitempty"`
	Product string             `form:"product,omitempty" json:"product,omitempty" bson:"product,omitempty" binding:"required"`
	Stock   int                `form:"stock" json:"stock" bson:"stock"`
	Color   string             `form:"color" json:"color" bson:"color" binding:"required"`
	Size    Size               `form:"size" json:"size" bson:"size" binding:"required"`
	Image   string             `form:"image" json:"image" bson:"image" binding:"required"`
}

type Size struct {
	Width  float64 `form:"width" json:"width" bson:"width" binding:"required"`
	Long   float64 `form:"long" json:"long" bson:"long" binding:"required"`
	Height float64 `form:"height" json:"height" bson:"height" binding:"required"`
}
