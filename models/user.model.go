package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	Location    string `form:"location" json:"location" bson:"location"`
	Road        string `form:"road" json:"road" bson:"road"`
	SubDistrict string `form:"subDistrict" json:"subDistrict" bson:"subDistrict"`
	District    string `form:"district" json:"district" bson:"district"`
	Pincode     int    `form:"pincode" json:"pincode" bson:"pincode"`
}

type User struct {
	ID            primitive.ObjectID `form:"_id,omitempty" json:"_id,omitempty" bson:"_id,omitempty"`
	UserName      string             `form:"userName" json:"userName" bson:"userName"`
	Password      string             `form:"password" json:"password" bson:"password"`
	Title         string             `form:"title" json:"title" bson:"title"`
	FirstName     string             `form:"firstName" json:"firstName" bson:"firstName"`
	LastName      string             `form:"lastName" json:"lastName" bson:"lastName"`
	BirthDate     string             `form:"birthDate" json:"birthDate" bson:"birthDate"`
	Address       Address            `form:"address" json:"address" bson:"address"`
	Gender        string             `form:"gender" json:"gender" bson:"gender"`
	Email         string             `form:"email" json:"email" bson:"email"`
	Phone         string             `form:"phone" json:"phone" bson:"phone"`
	Create_At     time.Time          `form:"create_at" json:"create_at" bson:"create_at"`
	Profile_Image string             `form:"profile_image" json:"profile_image" bson:"profile_image"`
}
