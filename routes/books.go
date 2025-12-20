// Package books for routes
package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tallyRead.com/models"
)

func CreateBook(context *gin.Context) {
	var book models.Book

	if err := context.ShouldBindJSON(&book); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, exists := context.Get("userId")

	if !exists {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No userID exists"})
		return
	}

	book.UserID = userId.(int64)

	err := book.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "book entry created successfully"})
}

func GetAllBooks(context *gin.Context) {

	userId, _ := context.Get("userId")

	books, err := models.GetAllBooks(userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": books})
}
