package controllers

import (
	"NHS-backend/models"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductController struct {
	productcollection          *mongo.Collection
	productAttributecollection *mongo.Collection
	ctx                        context.Context
}

func InitProduct(productcollection *mongo.Collection, productAttributecollection *mongo.Collection, ctx context.Context) ProductController {
	return ProductController{
		productcollection:          productcollection,
		productAttributecollection: productAttributecollection,
		ctx:                        ctx,
	}
}

func validateProduct(product models.Product) error {
	if product.Name == "" {
		return errors.New("product name cannot be empty")
	}

	if product.Type == "" {
		return errors.New("product type cannot be empty")
	}

	if product.Description == "" {
		return errors.New("product description cannot be empty")
	}

	if product.Image == "" {
		return errors.New("product image cannot be empty")
	}
	return nil
}

func validateProductAttribute(productAttribute models.ProductAttribute) error {
	if productAttribute.Color == "" {
		return errors.New("product attribute color cannot be empty")
	}

	if productAttribute.Image == "" {
		return errors.New("product attribute image cannot be empty")
	}

	if productAttribute.Size.Width <= 0 {
		return errors.New("product attribute width must be greater than zero")
	}

	if productAttribute.Size.Long <= 0 {
		return errors.New("product attribute long must be greater than zero")
	}

	if productAttribute.Size.Height <= 0 {
		return errors.New("product attribute height must be greater than zero")
	}
	return nil
}

func (productController *ProductController) CreateProduct(ctx *gin.Context) {

	requestProduct, _, err := ctx.Request.FormFile("product")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer requestProduct.Close()

	var productDecode models.Product
	err = json.NewDecoder(requestProduct).Decode((&productDecode))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	err = validateProduct(productDecode)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	requestProductAttribute, _, err := ctx.Request.FormFile("productAttribute")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer requestProductAttribute.Close()

	var productAttributeDecode []models.ProductAttribute
	err = json.NewDecoder(requestProductAttribute).Decode((&productAttributeDecode))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	productType, _ := primitive.ObjectIDFromHex(productDecode.Type)
	product := bson.D{
		bson.E{Key: "name", Value: productDecode.Name},
		bson.E{Key: "type", Value: productType},
		bson.E{Key: "description", Value: productDecode.Description},
		bson.E{Key: "image", Value: productDecode.Image},
	}
	newProduct, err := productController.productcollection.InsertOne(productController.ctx, product)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	var productAttributeInterface []interface{}
	for _, productAttributes := range productAttributeDecode {
		err = validateProductAttribute(productAttributes)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		productAttribute := bson.D{
			bson.E{Key: "product", Value: newProduct.InsertedID},
			bson.E{Key: "stock", Value: productAttributes.Stock},
			bson.E{Key: "color", Value: productAttributes.Color},
			bson.E{Key: "size", Value: productAttributes.Size},
			bson.E{Key: "image", Value: productAttributes.Image},
		}
		productAttributeInterface = append(productAttributeInterface, productAttribute)
	}

	_, e := productController.productAttributecollection.InsertMany(productController.ctx, productAttributeInterface)
	if e != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": e.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
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
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.All(ctx, &products); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
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
	productId := ctx.Param("id")
	objectProductId, _ := primitive.ObjectIDFromHex(productId)
	requestProduct, _, err := ctx.Request.FormFile("product")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer requestProduct.Close()

	var productDecode models.Product
	err = json.NewDecoder(requestProduct).Decode((&productDecode))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	err = validateProduct(productDecode)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	productType, _ := primitive.ObjectIDFromHex(productDecode.Type)
	filter := bson.D{bson.E{Key: "_id", Value: objectProductId}}
	product := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "name", Value: productDecode.Name},
		bson.E{Key: "type", Value: productType},
		bson.E{Key: "description", Value: productDecode.Description},
		bson.E{Key: "image", Value: productDecode.Image},
	}}}
	updateProduct, err := productController.productcollection.UpdateOne(productController.ctx, filter, product)
	if updateProduct.MatchedCount != 1 {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "no matched document found for update"})
	}
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	requestProductAttribute, _, err := ctx.Request.FormFile("productAttribute")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer requestProductAttribute.Close()

	var productAttributeDecode []models.ProductAttribute
	err = json.NewDecoder(requestProductAttribute).Decode((&productAttributeDecode))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	for _, productAttributes := range productAttributeDecode {
		err = validateProductAttribute(productAttributes)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		if productAttributes.ID == primitive.NilObjectID {
			productAttributes.ID = primitive.NewObjectID()
		}
		filter := bson.D{bson.E{Key: "_id", Value: productAttributes.ID}}
		productAttribute := bson.D{bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "product", Value: objectProductId},
			bson.E{Key: "stock", Value: productAttributes.Stock},
			bson.E{Key: "color", Value: productAttributes.Color},
			bson.E{Key: "size", Value: productAttributes.Size},
			bson.E{Key: "image", Value: productAttributes.Image},
		}}}
		opts := options.Update().SetUpsert(true)

		_, err := productController.productAttributecollection.UpdateOne(productController.ctx, filter, productAttribute, opts)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
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
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
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
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.All(ctx, &products); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
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

func (productController *ProductController) GetWithProductAttribute(ctx *gin.Context) {
	var product bson.M
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cursor, err := productController.productcollection.Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: objectId}}}},
		bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "productAttribute"}, {Key: "localField", Value: "_id"}, {Key: "foreignField", Value: "product"}, {Key: "as", Value: "productAttribute"}}}},
		// bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$productAttribute"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
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

func (productController *ProductController) Test(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (productController *ProductController) RegisterProductRoutes(rg *gin.RouterGroup) {
	productroute := rg.Group("/product")
	productroute.POST("/create", UploadMultipleFiles, productController.CreateProduct)
	productroute.GET("/get/:id", productController.GetProduct)
	productroute.GET("/getAll", productController.GetAll)
	productroute.PUT("/update/:id", UploadAndRemoveMultipleFiles, productController.UpdateProduct)
	productroute.DELETE("/delete/:id/:fileDelete", RemoveFile, productController.DeleteProduct)
	productroute.POST("/getProductForDatatable", productController.GetProductForDatatable)
	productroute.GET("/getWithProductAttribute/:id", productController.GetWithProductAttribute)
	productroute.POST("/test", productController.Test)
}
