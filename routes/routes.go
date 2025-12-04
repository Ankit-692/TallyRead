// Package routes for Routing
package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	server.POST("/register", RegisterUser)
	server.POST("/login", Login)
}
