// Package middlewares for auth
package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"tallyRead.com/utils"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized No Token Found"})
		return
	}

	tokenString := strings.TrimPrefix(token, "Bearer ")

	userId, err := utils.VerifyToken(tokenString)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized", "error": err.Error()})
		return
	}

	context.Set("userId", userId)

	context.Next()
}
