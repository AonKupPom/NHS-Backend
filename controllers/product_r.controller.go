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

type ProductRentController struct {
	productRentcollection *mongo.Collection
	ctx                   context.Context
}

func InitProductRent(productRentcollection *mongo.Collection, ctx context.Context) ProductRentController {
	return ProductRentController{
		productRentcollection: productRentcollection,
		ctx:                   ctx,
	}
}

func (productRentController *ProductRentController) CreateProductRent(ctx *gin.Context) {
	var productRent models.Product_rent
	if err := ctx.ShouldBind(&productRent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	newProductRent := bson.D{
		bson.E{Key: "product", Value: productRent.Product},
		bson.E{Key: "price", Value: productRent.Price},
	}
	_, err := productRentController.productRentcollection.InsertOne(productRentController.ctx, newProductRent)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (productRentController *ProductRentController) GetProductRent(ctx *gin.Context) {
	var productRent models.Product_rent
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{bson.E{Key: "_id", Value: objectId}}
	err := productRentController.productRentcollection.FindOne(ctx, query).Decode(&productRent)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, productRent)
}

func (productRentController *ProductRentController) GetAll(ctx *gin.Context) {
	var productRents []bson.M

	cursor, err := productRentController.productRentcollection.Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "price", Value: bson.D{{Key: "$gt", Value: 9}}}}}},
		bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "products"}, {Key: "localField", Value: "product"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "product"}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$product"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
	})

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	if err := cursor.All(ctx, &productRents); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	cursor.Close(ctx)

	if len(productRents) == 0 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "documents not found"})
		return
	}

	ctx.JSON(http.StatusOK, productRents)
}

func (productRentController *ProductRentController) UpdateProductRent(ctx *gin.Context) {
	var productRent models.Product_rent
	if err := ctx.ShouldBind(&productRent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "product", Value: productRent.Product},
		bson.E{Key: "price", Value: productRent.Price},
	}}}
	result, err := productRentController.productRentcollection.UpdateOne(ctx, filter, update)
	if result.MatchedCount != 1 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "no matched document found for update"})
	}
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (productRentController *ProductRentController) DeleteProductRent(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	result, err := productRentController.productRentcollection.DeleteOne(ctx, filter)
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

func (productRentController *ProductRentController) GetLazyProductRent(ctx *gin.Context) {
	type RequestBody struct {
		Skip int `form:"skip" json:"skip" bson:"skip"`
	}

	var productRents []bson.M
	var requestBody RequestBody

	if err := ctx.ShouldBind(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	cursor, err := productRentController.productRentcollection.Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "products"}, {Key: "localField", Value: "product"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "product"}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$product"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		bson.D{{Key: "$skip", Value: requestBody.Skip}},
		bson.D{{Key: "$limit", Value: 6}},
	})

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	if err := cursor.All(ctx, &productRents); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	cursor.Close(ctx)

	if len(productRents) == 0 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "documents not found"})
		return
	}

	ctx.JSON(http.StatusOK, productRents)
}

func (productRentController *ProductRentController) RegisterProductRentRoutes(rg *gin.RouterGroup) {
	productRentroute := rg.Group("/productRent")
	productRentroute.POST("/create", UploadFile, productRentController.CreateProductRent)
	productRentroute.GET("/get/:id", productRentController.GetProductRent)
	productRentroute.GET("/getAll", productRentController.GetAll)
	productRentroute.PUT("/update/:id", UploadFile, productRentController.UpdateProductRent)
	productRentroute.DELETE("/delete/:id", productRentController.DeleteProductRent)
	productRentroute.POST("/getLazyProductRent", productRentController.GetLazyProductRent)
}
