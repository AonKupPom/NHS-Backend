package controllers

import (
	"NHS-backend/models"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRentController struct {
	productRentcollection *mongo.Collection
	productcollection     *mongo.Collection
	ctx                   context.Context
}

func InitProductRent(
	productRentcollection *mongo.Collection,
	productcollection *mongo.Collection,
	ctx context.Context,
) ProductRentController {
	return ProductRentController{
		productRentcollection: productRentcollection,
		productcollection:     productcollection,
		ctx:                   ctx,
	}
}

func validateProductRent(productRent models.ProductRent) error {
	if productRent.Product == "" {
		return errors.New("productRent product cannot be empty")
	}

	if productRent.ProductAttribute == "" {
		return errors.New("productRent productAttribute cannot be empty")
	}

	if productRent.Price <= 0 {
		return errors.New("productRent price cannot be empty")
	}

	return nil
}

func (productRentController *ProductRentController) CreateProductRent(ctx *gin.Context) {
	var productRent models.ProductRent
	if err := ctx.ShouldBind(&productRent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	product, _ := primitive.ObjectIDFromHex(productRent.Product)
	newProductRent := bson.D{
		bson.E{Key: "product", Value: product},
		bson.E{Key: "price", Value: productRent.Price},
	}
	_, err := productRentController.productRentcollection.InsertOne(productRentController.ctx, newProductRent)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (productRentController *ProductRentController) GetProductRent(ctx *gin.Context) {
	var productRent models.ProductRent
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{bson.E{Key: "_id", Value: objectId}}
	err := productRentController.productRentcollection.FindOne(ctx, query).Decode(&productRent)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
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
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.All(ctx, &productRents); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
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
	productId := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(productId)
	enableRent, _ := strconv.ParseBool(ctx.PostForm("enableRent"))

	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	product := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "enableRent", Value: enableRent},
	}}}
	_, err := productRentController.productcollection.UpdateOne(productRentController.ctx, filter, product)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	requestProductRent, _, err := ctx.Request.FormFile("productRent")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer requestProductRent.Close()

	var productRentDecode []models.ProductRent
	err = json.NewDecoder(requestProductRent).Decode((&productRentDecode))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	for _, productRents := range productRentDecode {
		err = validateProductRent(productRents)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		if productRents.ID == primitive.NilObjectID {
			productRents.ID = primitive.NewObjectID()
		}
		productObjectId, _ := primitive.ObjectIDFromHex(productRents.Product)
		productRentObjectId, _ := primitive.ObjectIDFromHex(productRents.ProductAttribute)
		filter := bson.D{bson.E{Key: "_id", Value: productRents.ID}}
		productRent := bson.D{bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "product", Value: productObjectId},
			bson.E{Key: "productAttribute", Value: productRentObjectId},
			bson.E{Key: "price", Value: productRents.Price},
		}}}
		opts := options.Update().SetUpsert(true)

		_, err := productRentController.productRentcollection.UpdateOne(productRentController.ctx, filter, productRent, opts)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (productRentController *ProductRentController) DeleteProductRent(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	result, err := productRentController.productRentcollection.DeleteOne(ctx, filter)
	if result.DeletedCount != 1 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "no match document found for delete"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
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
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "products"},
			{Key: "localField", Value: "product"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "product"},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$product"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		bson.D{{Key: "$match", Value: bson.M{"product.enableRent": true}}},
		bson.D{{Key: "$skip", Value: requestBody.Skip}},
		bson.D{{Key: "$limit", Value: 6}},
	})

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.All(ctx, &productRents); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	cursor.Close(ctx)

	if len(productRents) == 0 {
		ctx.JSON(http.StatusOK, nil)
		return
	}

	ctx.JSON(http.StatusOK, productRents)
}

func (productRentController *ProductRentController) GetProductWithRentData(ctx *gin.Context) {
	var product bson.M
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cursor, err := productRentController.productcollection.Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: objectId}}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "productAttribute"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "product"},
			{Key: "as", Value: "productAttribute"},
			{Key: "pipeline", Value: bson.A{
				bson.D{{Key: "$project", Value: bson.M{"product": 0}}},
				bson.D{{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "productRent"},
					{Key: "localField", Value: "_id"},
					{Key: "foreignField", Value: "productAttribute"},
					{Key: "as", Value: "productRent"},
					{Key: "pipeline", Value: bson.A{bson.D{{Key: "$project", Value: bson.M{"_id": 1, "price": 1}}}}},
				}}},
				bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$productRent"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
			}},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "productRent"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "product"},
			{Key: "as", Value: "productRent"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "productType"},
			{Key: "localField", Value: "type"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "type"},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$type"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
	})

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if cursor.Next(ctx) {
		err := cursor.Decode(&product)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, product)
}

func (productRentController *ProductRentController) RegisterProductRentRoutes(rg *gin.RouterGroup) {
	productRentroute := rg.Group("/productRent")
	productRentroute.POST("/create", UploadFile, productRentController.CreateProductRent)
	productRentroute.GET("/get/:id", productRentController.GetProductRent)
	productRentroute.GET("/getAll", productRentController.GetAll)
	productRentroute.PUT("/update/:id", UploadFile, productRentController.UpdateProductRent)
	productRentroute.DELETE("/delete/:id", productRentController.DeleteProductRent)
	productRentroute.POST("/getLazyProductRent", productRentController.GetLazyProductRent)
	productRentroute.GET("/getProductWithRentData/:id", productRentController.GetProductWithRentData)
}
