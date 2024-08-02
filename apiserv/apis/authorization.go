package apiserv

import (
	"context"
	"net/http"
	"time"

	"seclolrip/mchosting/apiserv/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type DBUser struct {
	Token    string `bson:"token"`
	UserId   string `bson:"userId"`
	Username string `bson:"username"`
	Role     string `bson:"role"`
}

func (env *Env) Authorization(c *gin.Context) {
	userId := c.GetHeader("UserId")
	userToken, _ := c.Cookie("Token")

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

	c.Request.Header.Add("Role", userData.Role)

	c.Next()
}

func (env *Env) Login_Authorization(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userRaw := env.DB.Collection("users").FindOne(ctx, bson.M{"username": username, "password": password})
	if userRaw.Err() != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Username/Password is incorrect")
		return
	}

	var userData DBUser
	if err := userRaw.Decode(&userData); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "An Unexpected Error Occurred")
		return
	}

	tokenString, err := utils.JwtCreateToken(userData.Username+"//"+userData.UserId, userData.Role)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Unable to authorize credentials")
		return
	}

	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	_, updErr := env.DB.Collection("users").UpdateOne(ctx1, bson.M{"userId": userData.UserId}, bson.M{"token": tokenString, "tokenExp": time.Now().Add(time.Minute * 30).Unix()})
	if updErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Unable to authorize credentials")
		return
	}
	c.SetCookie("token", tokenString, 1800, "/", "mchosting", true, true)

	if redirectParam := c.Query("redirect"); redirectParam != "" {
		c.Redirect(http.StatusSeeOther, "/"+redirectParam)
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func (env *Env) CheckIsOwner(c *gin.Context) {
	if c.Request.Header.Get("Role") != "Owner" {
		c.AbortWithStatusJSON(http.StatusForbidden, "Not Allowed")
		return
	}

	c.Next()
}
