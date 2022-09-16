package main

import (
	"NHS-backend/controllers"
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	err         error

	usercollection *mongo.Collection
	tentcollection *mongo.Collection

	authcontroller controllers.AuthController
	usercontroller controllers.UserController
	tentcontroller controllers.TentController
)

func init() {
	ctx = context.TODO()

	mongoconn := options.Client().ApplyURI("mongodb+srv://sasawat:sIz9xe4LZFs9bvct@mean-stack.cb9amll.mongodb.net/?retryWrites=true&w=majority")
	mongoclient, err = mongo.Connect(ctx, mongoconn)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("mongo connection established")

	usercollection = mongoclient.Database("NHS-Database").Collection("users")
	tentcollection = mongoclient.Database("NHS-Database").Collection("tents")

	authcontroller = controllers.InitAuth(usercollection, ctx)
	usercontroller = controllers.InitUser(usercollection, ctx)
	tentcontroller = controllers.InitTent(tentcollection, ctx)

	server = gin.Default()
}

func main() {
	defer mongoclient.Disconnect(ctx)

	basepath := server.Group("/api")
	authcontroller.RegisterAuthRoutes(basepath)
	usercontroller.RegisterUserRoutes(basepath)
	tentcontroller.RegisterTentRoutes(basepath)

	log.Fatal(server.Run(":9090"))
}
