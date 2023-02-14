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

	productType, _ := primitive.ObjectIDFromHex(product.Type)
	newProduct := bson.D{
		bson.E{Key: "name", Value: product.Name},
		bson.E{Key: "type", Value: productType},
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
		ctx.JSON(http.StatusOK, product)
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
	productType, _ := primitive.ObjectIDFromHex(product.Type)
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "name", Value: product.Name},
		bson.E{Key: "type", Value: productType},
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

func (productController *ProductController) GetProductForDatatable(ctx *gin.Context) {
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

	var products []bson.M
	var requestBody RequestBody

	if err := ctx.ShouldBind(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	searchValue := bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: requestBody.Search, Options: "i"}}}
	count := bson.A{bson.D{{Key: "$match", Value: bson.D{{Key: "$or", Value: bson.A{
		bson.D{{Key: "name", Value: searchValue}},
		bson.D{{Key: "type", Value: searchValue}},
		bson.D{{Key: "description", Value: searchValue}},
	}}}}}, bson.D{{Key: "$count", Value: "count"}}}
	result := bson.A{bson.D{{Key: "$match", Value: bson.D{{Key: "$or", Value: bson.A{
		bson.D{{Key: "name", Value: searchValue}},
		bson.D{{Key: "type", Value: searchValue}},
		bson.D{{Key: "description", Value: searchValue}},
	}}}}},
		bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "productType"}, {Key: "localField", Value: "type"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "type"}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$type"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		bson.D{{Key: "$skip", Value: requestBody.Start}},
		bson.D{{Key: "$limit", Value: requestBody.TableRange}}}
	facetStage := bson.D{{Key: "$facet", Value: bson.D{{Key: "count", Value: count}, {Key: "result", Value: result}}}}

	cursor, err := productController.productcollection.Aggregate(ctx, mongo.Pipeline{facetStage, bson.D{{Key: "$unwind", Value: "$count"}}})

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	if err := cursor.All(ctx, &products); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	cursor.Close(ctx)

	emptyData := make([]string, 0)
	if len(products) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"recordsFiltered": 0,
			"recordsTotal":    0,
			"data":            emptyData,
		})
		return
	}

	var bsonToStruct BsonToStruct
	bsonBytes, _ := bson.Marshal(products[0])
	bson.Unmarshal(bsonBytes, &bsonToStruct)

	var response = Response{bsonToStruct.Count.Count, bsonToStruct.Count.Count, bsonToStruct.Result}

	ctx.JSON(http.StatusOK, response)
}

func (productController *ProductController) RegisterProductRoutes(rg *gin.RouterGroup) {
	productroute := rg.Group("/product")
	productroute.POST("/create", UploadFile, productController.CreateProduct)
	productroute.GET("/get/:id", productController.GetProduct)
	productroute.GET("/getAll", productController.GetAll)
	productroute.PUT("/update/:id", UploadAndRemoveFile, productController.UpdateProduct)
	productroute.DELETE("/delete/:id/:fileDelete", RemoveFile, productController.DeleteProduct)
	productroute.POST("/getProductForDatatable", productController.GetProductForDatatable)
}
