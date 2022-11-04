package controllers

import (
	"NHS-backend/models"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	var user models.User
	var requestBody models.LoginRequestBody

	if err := ctx.ShouldBind(&requestBody); err != nil {
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

func (authController *AuthController) GoogleLogin(ctx *gin.Context) {
	var user models.User
	var requestBody models.GoogleLoginRequestBody

	if err := ctx.ShouldBind(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	query := bson.D{bson.E{Key: "email", Value: requestBody.Email}}
	err := authController.usercollection.FindOne(ctx, query).Decode(&user)
	if err != nil {
		newUser := bson.D{
			bson.E{Key: "firstName", Value: requestBody.FirstName},
			bson.E{Key: "lastName", Value: requestBody.LastName},
			bson.E{Key: "email", Value: requestBody.Email},
			bson.E{Key: "create_at", Value: time.Now()},
		}
		insertUser, err := authController.usercollection.InsertOne(authController.ctx, newUser)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"meaasge": err.Error()})
			return
		}
		var userData models.User
		query := bson.D{bson.E{Key: "_id", Value: insertUser.InsertedID}}
		opts := options.FindOne().SetProjection(bson.D{{"userName", 0}, {"password", 0}})
		fatal := authController.usercollection.FindOne(ctx, query, opts).Decode(&userData)
		if fatal != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"meaasge": err.Error()})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"_id": userData.ID,
			// "exp": time.Now().Add(time.Minute * 1).Unix(),
		})
		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
		ctx.JSON(http.StatusOK, gin.H{"token": tokenString, "user": userData})
		return
	}
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
	authroute.POST("/googleLogin", authController.GoogleLogin)
}
