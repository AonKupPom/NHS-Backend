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

type ProductController struct {
	productcollection *mongo.Collection
	ctx               context.Context
}

func InitProduct(productcollection *mongo.Collection, ctx context.Context) ProductController {
	return ProductController{
		productcollection: productcollection,
		ctx:               ctx,
	}
}

func (productController *ProductController) CreateProduct(ctx *gin.Context) {
	var product models.Product
	if err := ctx.ShouldBind(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	newProduct := bson.D{
		bson.E{Key: "name", Value: product.Name},
		bson.E{Key: "type", Value: product.Type},
		bson.E{Key: "description", Value: product.Description},
		bson.E{Key: "image", Value: product.Image},
		bson.E{Key: "stock", Value: product.Stock},
	}
	_, err := productController.productcollection.InsertOne(productController.ctx, newProduct)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (productController *ProductController) GetProduct(ctx *gin.Context) {
	var product models.Product
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{bson.E{Key: "_id", Value: objectId}}
	err := productController.productcollection.FindOne(ctx, query).Decode(&product)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

func (productController *ProductController) GetAll(ctx *gin.Context) {
	var products []*models.Product
	cursor, err := productController.productcollection.Find(ctx, bson.D{{}})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	for cursor.Next(ctx) {
		var product models.Product
		err := cursor.Decode(&product)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			return
		}
		products = append(products, &product)
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	cursor.Close(ctx)

	if len(products) == 0 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "documents not found"})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (productController *ProductController) UpdateProduct(ctx *gin.Context) {
	var product models.Product
	if err := ctx.ShouldBind(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "name", Value: product.Name},
		bson.E{Key: "type", Value: product.Type},
		bson.E{Key: "description", Value: product.Description},
		bson.E{Key: "image", Value: product.Image},
		bson.E{Key: "stock", Value: product.Stock},
	}}}
	result, err := productController.productcollection.UpdateOne(ctx, filter, update)
	if result.MatchedCount != 1 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "no matched document found for update"})
	}
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (productController *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	result, err := productController.productcollection.DeleteOne(ctx, filter)
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

func (productController *ProductController) RegisterProductRoutes(rg *gin.RouterGroup) {
	productroute := rg.Group("/product")
	productroute.POST("/create", UploadFile, productController.CreateProduct)
	productroute.GET("/get/:id", productController.GetProduct)
	productroute.GET("/getAll", productController.GetAll)
	productroute.PUT("/update/:id", UploadFile, productController.UpdateProduct)
	productroute.DELETE("/delete/:id", productController.DeleteProduct)
}
