package main

import (
	"NHS-backend/middleware"
	"NHS-backend/router"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	err         error
)

func init() {
	ctx = context.TODO()

	mongoconn := options.Client().ApplyURI("mongodb+srv://sasawat:sIz9xe4LZFs9bvct@mean-stack.cb9amll.mongodb.net/?retryWrites=true&w=majority")
	mongoclient, err = mongo.Connect(ctx, mongoconn)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("mongo connection successfully")

	server = gin.Default()
	server.Use(middleware.CORSMiddleware())
	server.Use(static.Serve("/uploads", static.LocalFile("./uploads", true)))
}

func main() {

	defer mongoclient.Disconnect(ctx)

	router.InitRoute(server, mongoclient, ctx)

	log.Fatal(server.Run(os.Getenv("PORT")))
}
