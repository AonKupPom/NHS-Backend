package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bill struct {
	ID           primitive.ObjectID `form:"_id,omitempty" json:"_id,omitempty" bson:"_id,omitempty"`
	Payment_Type string             `form:"payment_type" json:"payment_type" bson:"payment_type" binding:"required"`
	Pay_Date     time.Time          `form:"pay_date" json:"pay_date" bson:"pay_date" binding:"required"`
	Total_Price  int                `form:"total_price" json:"total_price" bson:"total_price" binding:"required"`
	Bank         string             `form:"bank" json:"bank" bson:"bank" binding:"required"`
}
