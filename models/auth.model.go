package models

type LoginRequestBody struct {
	UserName string `form:"userName" json:"userName" bson:"userName"`
	Password string `form:"password" json:"password" bson:"password"`
}

type GoogleLoginRequestBody struct {
	Email     string `form:"email" json:"email" bson:"email"`
	FirstName string `form:"firstName" json:"firstName" bson:"firstName"`
	LastName  string `form:"lastName" json:"lastName" bson:"lastName"`
}
