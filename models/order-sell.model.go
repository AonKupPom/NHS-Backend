package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderSell struct {
	ID           primitive.ObjectID `form:"_id,omitempty" json:"_id,omitempty" bson:"_id,omitempty"`
	Product_Sell string             `form:"product_sell,omitempty" json:"product_sell,omitempty" bson:"product_sell,omitempty" binding:"required"`
	Customer     string             `form:"customer,omitempty" json:"customer,omitempty" bson:"customer,omitempty" binding:"required"`
	Buy_Date     time.Time          `form:"buy_date" json:"buy_date" bson:"buy_date" binding:"required"`
	Quantity     int                `form:"quantity" json:"quantity" bson:"quantity" binding:"required"`
	Bill         primitive.ObjectID `form:"bill" json:"bill" bson:"bill" binding:"required"`
}
