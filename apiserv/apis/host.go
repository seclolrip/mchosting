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

	authRoutes := router.Group("/")
	authRoutes.Use(env.Authorization)

	authRoutes.POST("/upload/mods", env.Upload_Mods)
	authRoutes.POST("/upload/modpacks", env.Upload_Modpacks)
	authRoutes.POST("/upload/resourcepacks", env.Upload_Resourcepacks)

	router.Run("localhost:8080")
}
