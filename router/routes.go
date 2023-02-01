package router

import (
	"NHS-backend/controllers"
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitRoute(server *gin.Engine, mongoclient *mongo.Client, ctx context.Context) {

	usercollection := mongoclient.Database("NHS-Database").Collection("users")
	productcollection := mongoclient.Database("NHS-Database").Collection("products")
	productSellcollection := mongoclient.Database("NHS-Database").Collection("productSell")
	productRentcollection := mongoclient.Database("NHS-Database").Collection("productRent")

	authcontroller := controllers.InitAuth(usercollection, ctx)
	usercontroller := controllers.InitUser(usercollection, ctx)
	productcontroller := controllers.InitProduct(productcollection, ctx)
	productSellcontroller := controllers.InitProductSell(productSellcollection, ctx)
	productRentcontroller := controllers.InitProductRent(productRentcollection, ctx)

	basepath := server.Group("/api")
	authcontroller.RegisterAuthRoutes(basepath)
	usercontroller.RegisterUserRoutes(basepath)
	productcontroller.RegisterProductRoutes(basepath)
	productSellcontroller.RegisterProductSellRoutes(basepath)
	productRentcontroller.RegisterProductRentRoutes(basepath)

}
