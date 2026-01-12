// Package middlewares for auth
package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tallyRead.com/utils"
)

func Authenticate(context *gin.Context) {
	token, err := context.Cookie("auth_token")

	if err != nil || token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized No Token Found"})
		return
	}

	userId, name, err := utils.VerifyToken(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized", "error": err.Error()})
		return
	}

	context.Set("userId", userId)
	context.Set("username", name)

	context.Next()
}

func Authorization(context *gin.Context) {
	username, _ := context.Get("username")
	context.JSON(http.StatusOK, gin.H{"authenticated": true, "username": username})
}
