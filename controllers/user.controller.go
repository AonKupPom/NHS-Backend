package controllers

import (
	"NHS-backend/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	usercollection *mongo.Collection
	ctx            context.Context
}

func InitUser(usercollection *mongo.Collection, ctx context.Context) UserController {
	return UserController{
		usercollection: usercollection,
		ctx:            ctx,
	}
}

func (userController *UserController) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	encryptPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	newUser := bson.D{
		bson.E{Key: "userName", Value: user.UserName},
		bson.E{Key: "password", Value: string(encryptPassword)},
		bson.E{Key: "title", Value: user.Title},
		bson.E{Key: "firstName", Value: user.FirstName},
		bson.E{Key: "lastName", Value: user.LastName},
		bson.E{Key: "birthDate", Value: user.BirthDate},
		bson.E{Key: "address", Value: user.Address},
		bson.E{Key: "gender", Value: user.Gender},
		bson.E{Key: "email", Value: user.Email},
		bson.E{Key: "phone", Value: user.Phone},
		bson.E{Key: "create_at", Value: user.Create_At},
	}

	_, err := userController.usercollection.InsertOne(userController.ctx, newUser)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"meaasge": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (userController *UserController) GetUser(ctx *gin.Context) {
	var user models.User
	opts := options.FindOne().SetProjection(bson.D{{"userName", 0}, {"password", 0}})
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{bson.E{Key: "_id", Value: objectId}}
	err := userController.usercollection.FindOne(ctx, query, opts).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"meaasge": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (userController *UserController) GetAll(ctx *gin.Context) {
	var users []*models.User
	cursor, err := userController.usercollection.Find(ctx, bson.D{{}})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	for cursor.Next(ctx) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			return
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	cursor.Close(ctx)

	if len(users) == 0 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "documents not found"})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (userController *UserController) UpdateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{bson.E{Key: "_id", Value: objectId}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "title", Value: user.Title},
		bson.E{Key: "firstName", Value: user.FirstName},
		bson.E{Key: "lastName", Value: user.LastName},
		bson.E{Key: "birthDate", Value: user.BirthDate},
		bson.E{Key: "address", Value: user.Address},
		bson.E{Key: "gender", Value: user.Gender},
		bson.E{Key: "email", Value: user.Email},
		bson.E{Key: "phone", Value: user.Phone},
		bson.E{Key: "create_at", Value: user.Create_At},
	}}}
	result, err := userController.usercollection.UpdateOne(ctx, query, update)
	if result.MatchedCount != 1 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "no matched document found for update"})
	}
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (userController *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{bson.E{Key: "_id", Value: objectId}}
	result, err := userController.usercollection.DeleteOne(ctx, query)
	if result.DeletedCount != 1 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "no match document foumd for dalete"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (userController *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	userroute := rg.Group("/user")
	userroute.POST("/register", userController.CreateUser)
	userroute.GET("/get/:id", userController.GetUser)
	userroute.GET("/getAll", userController.GetAll)
	userroute.PUT("/update/:id", userController.UpdateUser)
	userroute.DELETE("/delete/:id", userController.DeleteUser)
}
