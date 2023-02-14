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

type ProductTypeController struct {
	productTypecollection *mongo.Collection
	ctx                   context.Context
}

func InitProductType(productTypecollection *mongo.Collection, ctx context.Context) ProductTypeController {
	return ProductTypeController{
		productTypecollection: productTypecollection,
		ctx:                   ctx,
	}
}

func (productTypeController *ProductTypeController) CreateProductType(ctx *gin.Context) {
	var productType models.ProductType
	if err := ctx.ShouldBind(&productType); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	newProductType := bson.D{
		bson.E{Key: "name", Value: productType.Name},
		bson.E{Key: "description", Value: productType.Description},
	}
	_, err := productTypeController.productTypecollection.InsertOne(productTypeController.ctx, newProductType)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (productTypeController *ProductTypeController) GetProductType(ctx *gin.Context) {
	var productType models.ProductType
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{bson.E{Key: "_id", Value: objectId}}
	err := productTypeController.productTypecollection.FindOne(ctx, query).Decode(&productType)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, productType)
}

func (productTypeController *ProductTypeController) GetAll(ctx *gin.Context) {
	var productTypes []*models.ProductType
	cursor, err := productTypeController.productTypecollection.Find(ctx, bson.D{{}})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	for cursor.Next(ctx) {
		var productType models.ProductType
		err := cursor.Decode(&productType)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			return
		}
		productTypes = append(productTypes, &productType)
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	cursor.Close(ctx)

	if len(productTypes) == 0 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "documents not found"})
		return
	}

	ctx.JSON(http.StatusOK, productTypes)
}

func (productTypeController *ProductTypeController) UpdateProductType(ctx *gin.Context) {
	var productType models.ProductType
	if err := ctx.ShouldBind(&productType); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "name", Value: productType.Name},
		bson.E{Key: "description", Value: productType.Description},
	}}}
	result, err := productTypeController.productTypecollection.UpdateOne(ctx, filter, update)
	if result.MatchedCount != 1 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "no matched document found for update"})
	}
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (productTypeController *ProductTypeController) DeleteProductType(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	result, err := productTypeController.productTypecollection.DeleteOne(ctx, filter)
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

func (productTypeController *ProductTypeController) RegisterProductTypeRoutes(rg *gin.RouterGroup) {
	productTypeRoute := rg.Group("/productType")
	productTypeRoute.POST("/create", UploadFile, productTypeController.CreateProductType)
	productTypeRoute.GET("/get/:id", productTypeController.GetProductType)
	productTypeRoute.GET("/getAll", productTypeController.GetAll)
	productTypeRoute.PUT("/update/:id", UploadFile, productTypeController.UpdateProductType)
	productTypeRoute.DELETE("/delete/:id", productTypeController.DeleteProductType)
}
