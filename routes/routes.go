// Package routes for Routing
package routes

import (
	"github.com/gin-gonic/gin"
	"tallyRead.com/middlewares"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/register", RegisterUser)
	server.POST("/login", Login)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.POST("/addBook", CreateBook)
	authenticated.GET("/getAllBooks", GetAllBooks)
}
