package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderRent struct {
	ID           primitive.ObjectID `form:"_id,omitempty" json:"_id,omitempty" bson:"_id,omitempty"`
	Product_Rent string             `form:"product_rent,omitempty" json:"product_rent,omitempty" bson:"product_rent,omitempty" binding:"required"`
	Customer     string             `form:"customer,omitempty" json:"customer,omitempty" bson:"customer,omitempty" binding:"required"`
	Rent_Date    time.Time          `form:"rent_date" json:"rent_date" bson:"rent_date" binding:"required"`
	Return_Date  time.Time          `form:"return_date" json:"return_date" bson:"return_date" binding:"required"`
	Rent_Status  bool               `form:"rent_status" json:"rent_status" bson:"rent_status" binding:"required"`
	Quantity     int                `form:"quantity" json:"quantity" bson:"quantity" binding:"required"`
}
