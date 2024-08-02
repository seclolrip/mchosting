package apiserv

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Env struct {
	DB *mongo.Database
}

func Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("srv"))
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	env := &Env{client.Database("databaseName")}

	router := gin.Default()

	router.POST("/login", env.Login_Authorization)

	createRoutes := router.Group("/create")
	createRoutes.Use(env.Authorization)

	createRoutes.POST("/server", env.CreatePendingServer)

	uploadRoutes := router.Group("/upload")
	uploadRoutes.Use(env.Authorization, env.CheckIsOwner)

	uploadRoutes.POST("/mods", env.Upload_Mods)
	uploadRoutes.POST("/modpacks", env.Upload_Modpacks)
	uploadRoutes.POST("/resourcepacks", env.Upload_Resourcepacks)

	router.Run("localhost:8080")
}
