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

type ProductSellController struct {
	productSellcollection *mongo.Collection
	ctx                   context.Context
}

func InitProductSell(productSellcollection *mongo.Collection, ctx context.Context) ProductSellController {
	return ProductSellController{
		productSellcollection: productSellcollection,
		ctx:                   ctx,
	}
}

func (productSellController *ProductSellController) CreateProductSell(ctx *gin.Context) {
	var productSell models.Product_sell
	if err := ctx.ShouldBind(&productSell); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	newProductSell := bson.D{
		bson.E{Key: "product", Value: productSell.Product},
		bson.E{Key: "price", Value: productSell.Price},
	}
	_, err := productSellController.productSellcollection.InsertOne(productSellController.ctx, newProductSell)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (productSellController *ProductSellController) GetProductSell(ctx *gin.Context) {
	var productSell models.Product_sell
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{bson.E{Key: "_id", Value: objectId}}
	err := productSellController.productSellcollection.FindOne(ctx, query).Decode(&productSell)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, productSell)
}

func (productSellController *ProductSellController) GetAll(ctx *gin.Context) {
	var productSells []*models.Product_sell
	cursor, err := productSellController.productSellcollection.Find(ctx, bson.D{{}})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	for cursor.Next(ctx) {
		var productSell models.Product_sell
		err := cursor.Decode(&productSell)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			return
		}
		productSells = append(productSells, &productSell)
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	cursor.Close(ctx)

	if len(productSells) == 0 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "documents not found"})
		return
	}

	ctx.JSON(http.StatusOK, productSells)
}

func (productSellController *ProductSellController) UpdateProductSell(ctx *gin.Context) {
	var productSell models.Product_sell
	if err := ctx.ShouldBind(&productSell); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "product", Value: productSell.Product},
		bson.E{Key: "price", Value: productSell.Price},
	}}}
	result, err := productSellController.productSellcollection.UpdateOne(ctx, filter, update)
	if result.MatchedCount != 1 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "no matched document found for update"})
	}
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (productSellController *ProductSellController) DeleteProductSell(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	result, err := productSellController.productSellcollection.DeleteOne(ctx, filter)
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

func (productSellController *ProductSellController) RegisterProductSellRoutes(rg *gin.RouterGroup) {
	productSellroute := rg.Group("/productSell")
	productSellroute.POST("/create", UploadFile, productSellController.CreateProductSell)
	productSellroute.GET("/get/:id", productSellController.GetProductSell)
	productSellroute.GET("/getAll", productSellController.GetAll)
	productSellroute.PUT("/update/:id", UploadFile, productSellController.UpdateProductSell)
	productSellroute.DELETE("/delete/:id", productSellController.DeleteProductSell)
}
