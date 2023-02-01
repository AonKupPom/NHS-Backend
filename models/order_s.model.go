package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order_sell struct {
	ID           primitive.ObjectID `form:"_id,omitempty" json:"_id,omitempty" bson:"_id,omitempty"`
	Product_Sell primitive.ObjectID `form:"product_sell,omitempty" json:"product_sell,omitempty" bson:"product_sell,omitempty"`
	Customer     primitive.ObjectID `form:"customer,omitempty" json:"customer,omitempty" bson:"customer,omitempty"`
	Buy_Date     time.Time          `form:"buy_date" json:"buy_date" bson:"buy_date"`
	Quantity     int                `form:"quantity" json:"quantity" bson:"quantity"`
	Bill         primitive.ObjectID `form:"bill" json:"bill" bson:"bill"`
}
