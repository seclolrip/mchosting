package apiserv

import (
	"fmt"
	"net/http"

	"seclolrip/mchosting/apiserv/utils"

	"github.com/gin-gonic/gin"
)

func (env *Env) Login_Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	tokenString, err := utils.JwtCreateToken(username)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating token")
		return
	}

	fmt.Printf("Token created: %s\n", tokenString)
	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusSeeOther, "/")

	c.String(http.StatusUnauthorized, "Invalid credentials")
}
