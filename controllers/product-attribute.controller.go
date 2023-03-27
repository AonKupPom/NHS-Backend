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

type ProductAttributeController struct {
	productAttributecollection *mongo.Collection
	ctx                        context.Context
}

func InitProductAttribute(productAttributecollection *mongo.Collection, ctx context.Context) ProductAttributeController {
	return ProductAttributeController{
		productAttributecollection: productAttributecollection,
		ctx:                        ctx,
	}
}

func (productAttributeController *ProductAttributeController) CreateProductAttribute(ctx *gin.Context) {
	var productAttribute models.ProductAttribute
	if err := ctx.ShouldBind(&productAttribute); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	product, _ := primitive.ObjectIDFromHex(productAttribute.Product)
	newProductAttribute := bson.D{
		bson.E{Key: "product", Value: product},
		bson.E{Key: "stock", Value: productAttribute.Stock},
		bson.E{Key: "color", Value: productAttribute.Color},
		bson.E{Key: "size", Value: productAttribute.Size},
	}
	_, err := productAttributeController.productAttributecollection.InsertOne(productAttributeController.ctx, newProductAttribute)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (productAttributeController *ProductAttributeController) GetProductAttribute(ctx *gin.Context) {
	var productAttribute models.ProductAttribute
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{bson.E{Key: "_id", Value: objectId}}
	err := productAttributeController.productAttributecollection.FindOne(ctx, query).Decode(&productAttribute)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, productAttribute)
}

func (productAttributeController *ProductAttributeController) GetAll(ctx *gin.Context) {
	var productAttributes []*models.ProductAttribute
	cursor, err := productAttributeController.productAttributecollection.Find(ctx, bson.D{{}})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.All(ctx, &productAttributes); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	cursor.Close(ctx)

	if len(productAttributes) == 0 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "documents not found"})
		return
	}

	ctx.JSON(http.StatusOK, productAttributes)
}

func (productAttributeController *ProductAttributeController) UpdateProductAttribute(ctx *gin.Context) {
	var productAttribute models.ProductAttribute
	if err := ctx.ShouldBind(&productAttribute); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	product, _ := primitive.ObjectIDFromHex(productAttribute.Product)
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "product", Value: product},
		bson.E{Key: "stock", Value: productAttribute.Stock},
		bson.E{Key: "color", Value: productAttribute.Color},
		bson.E{Key: "size", Value: productAttribute.Size},
	}}}
	result, err := productAttributeController.productAttributecollection.UpdateOne(ctx, filter, update)
	if result.MatchedCount != 1 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "no matched document found for update"})
	}
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (productAttributeController *ProductAttributeController) DeleteProductAttribute(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	result, err := productAttributeController.productAttributecollection.DeleteOne(ctx, filter)
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

func (productAttributeController *ProductAttributeController) RegisterProductAttributeRoutes(rg *gin.RouterGroup) {
	productAttributeRoute := rg.Group("/productAttribute")
	productAttributeRoute.POST("/create", UploadFile, productAttributeController.CreateProductAttribute)
	productAttributeRoute.GET("/get/:id", productAttributeController.GetProductAttribute)
	productAttributeRoute.GET("/getAll", productAttributeController.GetAll)
	productAttributeRoute.PUT("/update/:id", UploadFile, productAttributeController.UpdateProductAttribute)
	productAttributeRoute.DELETE("/delete/:id", productAttributeController.DeleteProductAttribute)
}
