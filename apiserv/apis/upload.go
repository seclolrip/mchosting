package apiserv

import (
	"context"
	"net/http"
	"time"

	"seclolrip/mchosting/apiserv/queries"
	"seclolrip/mchosting/apiserv/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (env *Env) Upload_Mods(c *gin.Context) {
	userId := c.GetHeader("UserId")

	server, err := queries.GetPendingServer(env.DB, userId)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		env.DB.Collection("bugs").InsertOne(ctx, bson.M{"type": "api", "route": "/upload/mods", "error": err.Error()})
		c.String(http.StatusInternalServerError, "Failed to find server! A bug has been reported!")
		return
	}

	intServer, err := queries.GetInternalServer(env.DB, server.ServerUUID)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		env.DB.Collection("bugs").InsertOne(ctx, bson.M{"type": "api", "route": "/upload/mods", "error": err.Error()})
		c.String(http.StatusInternalServerError, "Failed to upload mods! A bug has been reported!")
		return
	}

	file, _ := c.FormFile("file")

	uploadErr := utils.UploadFileToInternal(file, intServer.PrivIP, intServer.UUID.String(), "mods.gzip")
	if uploadErr != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		env.DB.Collection("bugs").InsertOne(ctx, bson.M{"type": "api", "route": "/upload/mods", "error": uploadErr.Error()})
		c.String(http.StatusInternalServerError, "Failed to upload mods! A bug has been reported!")
		return
	}

	c.String(http.StatusOK, "Uploaded")
}

func (env *Env) Upload_Modpacks(c *gin.Context) {
	userId := c.GetHeader("UserId")

	server, err := queries.GetPendingServer(env.DB, userId)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		env.DB.Collection("bugs").InsertOne(ctx, bson.M{"type": "api", "route": "/upload/modpacks", "error": err.Error()})
		c.String(http.StatusInternalServerError, "Failed to find server! A bug has been reported!")
		return
	}

	intServer, err := queries.GetInternalServer(env.DB, server.ServerUUID)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		env.DB.Collection("bugs").InsertOne(ctx, bson.M{"type": "api", "route": "/upload/modpacks", "error": err.Error()})
		c.String(http.StatusInternalServerError, "Failed to upload modpacks! A bug has been reported!")
		return
	}

	file, _ := c.FormFile("file")

	uploadErr := utils.UploadFileToInternal(file, intServer.PrivIP, intServer.UUID.String(), "modpack.gzip")
	if uploadErr != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		env.DB.Collection("bugs").InsertOne(ctx, bson.M{"type": "api", "route": "/upload/modpacks", "error": uploadErr.Error()})
		c.String(http.StatusInternalServerError, "Failed to upload modpacks! A bug has been reported!")
		return
	}

	c.String(http.StatusOK, "Uploaded")
}

func (env *Env) Upload_Resourcepacks(c *gin.Context) {
	userId := c.GetHeader("UserId")

	server, err := queries.GetPendingServer(env.DB, userId)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		env.DB.Collection("bugs").InsertOne(ctx, bson.M{"type": "api", "route": "/upload/resourcepacks", "error": err.Error()})
		c.String(http.StatusInternalServerError, "Failed to find server! A bug has been reported!")
		return
	}

	intServer, err := queries.GetInternalServer(env.DB, server.ServerUUID)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		env.DB.Collection("bugs").InsertOne(ctx, bson.M{"type": "api", "route": "/upload/resourcepacks", "error": err.Error()})
		c.String(http.StatusInternalServerError, "Failed to upload resourcepacks! A bug has been reported!")
		return
	}

	file, _ := c.FormFile("file")

	uploadErr := utils.UploadFileToInternal(file, intServer.PrivIP, intServer.UUID.String(), "resourcepacks.gzip")
	if uploadErr != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		env.DB.Collection("bugs").InsertOne(ctx, bson.M{"type": "api", "route": "/upload/resourcepacks", "error": uploadErr.Error()})
		c.String(http.StatusInternalServerError, "Failed to upload resourcepacks! A bug has been reported!")
		return
	}

	c.String(http.StatusOK, "Uploaded")
}
