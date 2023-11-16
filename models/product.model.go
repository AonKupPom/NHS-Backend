package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `form:"_id,omitempty" json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `form:"name" json:"name" bson:"name" binding:"required"`
	Type        string             `form:"type,omitempty" json:"type,omitempty" bson:"type,omitempty" binding:"required"`
	Description string             `form:"description" json:"description" bson:"description" binding:"required"`
	Image       string             `form:"image" json:"image" bson:"image" binding:"required"`
	EnableRent  bool               `form:"enableRent" json:"enableRent" bson:"enableRent" binding:"required"`
	EnableSell  bool               `form:"enableSell" json:"enableSell" bson:"enableSell" binding:"required"`
}
