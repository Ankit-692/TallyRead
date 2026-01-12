// Package routes for Routing
package routes

import (
	"github.com/gin-gonic/gin"
	"tallyRead.com/middlewares"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/register", RegisterUser)
	server.POST("/login", Login)
	server.POST("/forgot-Password", ForgotPassword)
	server.POST("/reset-password", ResetPassword)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.POST("/addBook", CreateBook)
	authenticated.GET("/getAllBooks", GetAllBooks)
	authenticated.DELETE("/deleteBook/:id", DeleteBookByID)
	authenticated.GET("/book/:id", GetBookByID)
	authenticated.POST("/book/:id", UpdateBook)
	authenticated.GET("/api/me", middlewares.Authorization)
	authenticated.POST("/logout", Logout)
}
