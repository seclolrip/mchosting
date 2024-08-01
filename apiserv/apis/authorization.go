package apiserv

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type DBUser struct {
	Token  string
	UserId string
}

func (env *Env) Authorization(c *gin.Context) {
	userId := c.GetHeader("userId")
	userToken := c.GetHeader("Token")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userRaw := env.DB.Collection("users").FindOne(ctx, bson.M{"userId": userId})
	if userRaw.Err() != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	var userData DBUser
	if err := userRaw.Decode(&userData); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "An Unexpected Error Occurred")
		return
	} else if userData.Token != userToken {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	c.Next()
}
