package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order_rent struct {
	ID           primitive.ObjectID `form:"_id,omitempty" json:"_id,omitempty" bson:"_id,omitempty"`
	Product_Rent primitive.ObjectID `form:"product_rent,omitempty" json:"product_rent,omitempty" bson:"product_rent,omitempty"`
	Customer     primitive.ObjectID `form:"customer,omitempty" json:"customer,omitempty" bson:"customer,omitempty"`
	Rent_Date    time.Time          `form:"rent_date" json:"rent_date" bson:"rent_date"`
	Return_Date  time.Time          `form:"return_date" json:"return_date" bson:"return_date"`
	Rent_Status  bool               `form:"rent_status" json:"rent_status" bson:"rent_status"`
	Quantity     int                `form:"quantity" json:"quantity" bson:"quantity"`
}
