package controllers

import (
	"NHS-backend/models"
	"context"
	"net/http"

	"time"

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
	if err := ctx.ShouldBind(&user); err != nil {
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
		bson.E{Key: "create_at", Value: time.Now()},
		bson.E{Key: "role", Value: user.Role},
	}

	_, err := userController.usercollection.InsertOne(userController.ctx, newUser)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (userController *UserController) GetUser(ctx *gin.Context) {
	var user models.User
	opts := options.FindOne().SetProjection(bson.D{{Key: "userName", Value: 0}, {Key: "password", Value: 0}})
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{bson.E{Key: "_id", Value: objectId}}
	err := userController.usercollection.FindOne(ctx, query, opts).Decode(&user) //Decode ใช้เพื่อแปลง cursor ให้เป็น bson object
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			ctx.JSON(http.StatusOK, nil)
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (userController *UserController) GetAll(ctx *gin.Context) {
	var users []*models.User
	opts := options.Find().SetProjection(bson.D{{Key: "userName", Value: 0}, {Key: "password", Value: 0}})
	cursor, err := userController.usercollection.Find(ctx, bson.D{{}}, opts)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.All(ctx, &users); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
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
	if err := ctx.ShouldBind(&user); err != nil {
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
		bson.E{Key: "role", Value: user.Role},
	}}}
	result, err := userController.usercollection.UpdateOne(ctx, query, update)
	if result.MatchedCount != 1 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "no matched document found for update"})
	}
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
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
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (userController *UserController) GetUserForDatatable(ctx *gin.Context) {
	type RequestBody struct {
		Start      int    `form:"start" json:"start" bson:"start"`
		TableRange int    `form:"tableRange" json:"tableRange" bson:"tableRange"`
		Search     string `form:"search" json:"search" bson:"search"`
	}

	type counts struct {
		Count int `form:"count" json:"count" bson:"count"`
	}

	type BsonToStruct struct {
		Count  counts   `form:"count" json:"count" bson:"count"`
		Result []bson.M `form:"result" json:"result" bson:"result"`
	}

	type Response struct {
		RecordsFiltered int      `form:"recordsFiltered" json:"recordsFiltered" bson:"recordsFiltered"`
		RecordsTotal    int      `form:"recordsTotal" json:"recordsTotal" bson:"recordsTotal"`
		Data            []bson.M `form:"data" json:"data" bson:"data"`
	}

	var users []bson.M
	var requestBody RequestBody

	if err := ctx.ShouldBind(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	searchValue := bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: requestBody.Search, Options: "i"}}}
	count := bson.A{bson.D{{Key: "$match", Value: bson.D{{Key: "$or", Value: bson.A{
		bson.D{{Key: "firstName", Value: searchValue}},
		bson.D{{Key: "lastName", Value: searchValue}},
	}}}}}, bson.D{{Key: "$count", Value: "count"}}}
	result := bson.A{bson.D{{Key: "$match", Value: bson.D{{Key: "$or", Value: bson.A{
		bson.D{{Key: "firstName", Value: searchValue}},
		bson.D{{Key: "lastName", Value: searchValue}},
	}}}}},
		bson.D{{Key: "$skip", Value: requestBody.Start}},
		bson.D{{Key: "$limit", Value: requestBody.TableRange}}}
	facetStage := bson.D{{Key: "$facet", Value: bson.D{{Key: "count", Value: count}, {Key: "result", Value: result}}}}

	cursor, err := userController.usercollection.Aggregate(ctx, mongo.Pipeline{facetStage, bson.D{{Key: "$unwind", Value: "$count"}}})

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.All(ctx, &users); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	cursor.Close(ctx)

	emptyData := make([]string, 0)
	if len(users) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"recordsFiltered": 0,
			"recordsTotal":    0,
			"data":            emptyData,
		})
		return
	}

	var bsonToStruct BsonToStruct
	bsonBytes, _ := bson.Marshal(users[0])
	bson.Unmarshal(bsonBytes, &bsonToStruct)

	var response = Response{bsonToStruct.Count.Count, bsonToStruct.Count.Count, bsonToStruct.Result}

	ctx.JSON(http.StatusOK, response)
}

func (userController *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	userroute := rg.Group("/user")
	userroute.POST("/register", userController.CreateUser)
	userroute.GET("/get/:id", userController.GetUser)
	userroute.GET("/getAll", userController.GetAll)
	userroute.PUT("/update/:id", userController.UpdateUser)
	userroute.DELETE("/delete/:id", userController.DeleteUser)
	userroute.POST("/getUserForDatatable", userController.GetUserForDatatable)
}
