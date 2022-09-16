package controllers

import (
	"NHS-backend/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TentController struct {
	tentcollection *mongo.Collection
	ctx            context.Context
}

func InitTent(tentcollection *mongo.Collection, ctx context.Context) TentController {
	return TentController{
		tentcollection: tentcollection,
		ctx:            ctx,
	}
}

func (tentController *TentController) CreateTent(ctx *gin.Context) {
	var tent models.Tent
	if err := ctx.ShouldBindJSON(&tent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	_, err := tentController.tentcollection.InsertOne(tentController.ctx, tent)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"meaasge": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (tentController *TentController) GetTent(ctx *gin.Context) {
	var tent models.Tent
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{bson.E{Key: "_id", Value: objectId}}
	err := tentController.tentcollection.FindOne(ctx, query).Decode(&tent)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"meaasge": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tent)
}

func (tentController *TentController) GetAll(ctx *gin.Context) {
	var tents []*models.Tent
	cursor, err := tentController.tentcollection.Find(ctx, bson.D{{}})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	for cursor.Next(ctx) {
		var tent models.Tent
		err := cursor.Decode(&tent)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			return
		}
		tents = append(tents, &tent)
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	cursor.Close(ctx)

	if len(tents) == 0 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "documents not found"})
		return
	}

	ctx.JSON(http.StatusOK, tents)
}

func (tentController *TentController) UpdateTent(ctx *gin.Context) {
	var tent models.Tent
	if err := ctx.ShouldBindJSON(&tent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "name", Value: tent.Name},
		bson.E{Key: "size", Value: tent.Size},
		bson.E{Key: "color", Value: tent.Color},
		bson.E{Key: "price", Value: tent.Price},
		bson.E{Key: "type", Value: tent.Type},
		bson.E{Key: "shape", Value: tent.Shape},
		bson.E{Key: "description", Value: tent.Description},
		bson.E{Key: "image", Value: tent.Image},
		bson.E{Key: "stock", Value: tent.Stock},
	}}}
	result, err := tentController.tentcollection.UpdateOne(ctx, filter, update)
	if result.MatchedCount != 1 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "no matched document found for update"})
	}
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (tentController *TentController) DeleteTent(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	result, err := tentController.tentcollection.DeleteOne(ctx, filter)
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

func (tentController *TentController) RegisterTentRoutes(rg *gin.RouterGroup) {
	tentroute := rg.Group("/tent")
	tentroute.POST("/create", tentController.CreateTent)
	tentroute.GET("/get/:id", tentController.GetTent)
	tentroute.GET("/getAll", tentController.GetAll)
	tentroute.PUT("/update/:id", tentController.UpdateTent)
	tentroute.DELETE("/delete/:id", tentController.DeleteTent)
}
