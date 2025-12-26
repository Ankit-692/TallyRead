// Package books for routes
package routes

import (
	"net/http"
	"strconv"
	"strings"

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
		if strings.Contains(err.Error(), "2067") {
			context.JSON(http.StatusConflict, gin.H{"error": "Entry Already Exists"})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "book entry created successfully"})
}

func GetAllBooks(context *gin.Context) {

	userId, _ := context.Get("userId")

	status := context.DefaultQuery("status", "all")

	if status == "all" {
		books, err := models.GetAllBooks(userId)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"books": books})
	} else {
		books, err := models.GetBooksByStatus(userId.(int64), status)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"books": books})
	}
}

func DeleteBookByID(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := models.GetBookByID(id)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = book.Delete()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func GetBookByID(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := models.GetBookByID(id)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"book": book})

}

func UpdateBook(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedBook models.Book

	if err = context.ShouldBindJSON(&updatedBook); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	book, err := models.GetBookByID(id)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if book.TotalPage < updatedBook.PageRead {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Number of pages read can not be more than total Pages"})
		return
	}

	book.PageRead = updatedBook.PageRead
	book.State = updatedBook.State

	err = book.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Book Updated successfully"})

}
