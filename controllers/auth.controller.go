package controllers

import (
	"NHS-backend/models"
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	usercollection *mongo.Collection
	ctx            context.Context
}

func InitAuth(usercollection *mongo.Collection, ctx context.Context) AuthController {
	godotenv.Load(".env")

	return AuthController{
		usercollection: usercollection,
		ctx:            ctx,
	}
}

func (authController *AuthController) Login(ctx *gin.Context) {

	type LoginRequestBody struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}

	var user models.User
	var requestBody LoginRequestBody

	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	query := bson.D{bson.E{Key: "userName", Value: requestBody.UserName}}
	err := authController.usercollection.FindOne(ctx, query).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"meaasge": "ไม่มีชื่อผู้ใช้นี้ในระบบ กรุณาตรวจสอบอีกครั้ง"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"meaasge": "รหัสผ่านไม่ถูกต้อง กรุณาตรวจสอบอีกครั้ง"})
		return
	}

	user.UserName = ""
	user.Password = ""

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id": user.ID,
		// "exp": time.Now().Add(time.Minute * 1).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString, "user": user})
}

func (authController *AuthController) RegisterAuthRoutes(rg *gin.RouterGroup) {
	authroute := rg.Group("/auth")
	authroute.POST("/login", authController.Login)
}
