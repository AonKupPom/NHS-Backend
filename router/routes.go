package router

import (
	"NHS-backend/controllers"
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitRoute(server *gin.Engine, mongoclient *mongo.Client, ctx context.Context) {

	usercollection := mongoclient.Database("NHS-Database").Collection("users")
	tentcollection := mongoclient.Database("NHS-Database").Collection("tents")

	authcontroller := controllers.InitAuth(usercollection, ctx)
	usercontroller := controllers.InitUser(usercollection, ctx)
	tentcontroller := controllers.InitTent(tentcollection, ctx)

	basepath := server.Group("/api")
	authcontroller.RegisterAuthRoutes(basepath)
	usercontroller.RegisterUserRoutes(basepath)
	tentcontroller.RegisterTentRoutes(basepath)
}
