// Package books for routes
package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"tallyRead.com/db"
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

func SearchBooks(context *gin.Context) {
	query := context.Query("q")
	cacheKey := "books:" + query

	// 1. Check Redis
	if cachedData, err := db.RedisClient.Get(context, cacheKey).Result(); err == nil {
		context.Data(http.StatusOK, "application/json", []byte(cachedData))
		return
	}

	// 2. Fetch from Google
	apiKey := os.Getenv("GOOGLE_BOOKS_API_KEY")
	client := resty.New()
	var googleResp struct {
		Items []struct {
			VolumeInfo struct {
				Title         string   `json:"title"`
				Description   string   `json:"description"`
				Authors       []string `json:"authors"`
				PageCount     int16    `json:"pageCount"`
				AverageRating float64  `json:"averageRating"`
				ImageLinks    struct {
					Thumbnail string `json:"thumbnail"`
				} `json:"imageLinks"`
				PublishedDate string `json:"publishedDate"`
			} `json:"volumeInfo"`
		} `json:"items"`
	}

	_, err := client.R().SetQueryParams(map[string]string{
		"q": query, "maxResults": "40", "key": apiKey,
	}).SetResult(&googleResp).Get("https://www.googleapis.com/books/v1/volumes")

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "API failed"})
		return
	}

	// 3. Map to your Book Struct
	var results []models.Book
	for _, item := range googleResp.Items {
		info := item.VolumeInfo
		authors := "Unknown Author"
		if len(info.Authors) > 0 {
			authors = strings.Join(info.Authors, ", ")
		}

		img := info.ImageLinks.Thumbnail
		if img == "" {
			img = "assets/noCover.png"
		}

		results = append(results, models.Book{
			Title:         info.Title,
			Description:   info.Description,
			Authors:       authors,
			TotalPage:     info.PageCount,
			Ratings:       fmt.Sprintf("%.1f", info.AverageRating),
			Image:         img,
			PublishedDate: info.PublishedDate,
		})
	}

	// 4. Cache the clean Slice
	jsonData, _ := json.Marshal(results)
	db.RedisClient.Set(context, cacheKey, jsonData, 24*time.Hour)

	context.JSON(http.StatusOK, results)
}
