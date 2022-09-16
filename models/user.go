package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	Location    string `json:"location" bson:"location"`
	Road        string `json:"road" bson:"road"`
	SubDistrict string `json:"subDistrict" bson:"subDistrict"`
	District    string `json:"district" bson:"district"`
	Pincode     int    `json:"pincode" bson:"pincode"`
}

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName  string             `json:"userName" bson:"userName"`
	Password  string             `json:"password" bson:"password"`
	Title     string             `json:"title" bson:"title"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	BirthDate string             `json:"birthDate" bson:"birthDate"`
	Address   Address            `json:"address" bson:"address"`
	Gender    string             `json:"gender" bson:"gender"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"Phone" bson:"phone"`
	Create_At string             `json:"create_at" bson:"create_at"`
}
