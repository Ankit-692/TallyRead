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

	authenticated := server.Group("/api")
	authenticated.Use(middlewares.Authenticate)
	authenticated.POST("/addBook", CreateBook)
	authenticated.GET("/getAllBooks", GetAllBooks)
	authenticated.POST("/deleteBook/:id", DeleteBookByID)
	authenticated.GET("/book/:id", GetBookByID)
	authenticated.POST("/book/:id", UpdateBook)
	authenticated.GET("/me", middlewares.Authorization)
	authenticated.GET("/searchBooks", SearchBooks)
	authenticated.POST("/logout", Logout)
}
