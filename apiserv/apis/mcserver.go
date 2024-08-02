package apiserv

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (env *Env) CreatePendingServer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pendingServer, pendingErr := env.DB.Collection("pendingServers").InsertOne(ctx, bson.M{"userId": c.GetHeader("UserId"), "name": ""})
	if pendingErr != nil {
		ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel1()

		env.DB.Collection("pendingServers").InsertOne(ctx1, bson.M{"type": "api", "route": "/upload/mods", "error": pendingErr.Error()})
	}
}
