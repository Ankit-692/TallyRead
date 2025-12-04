// Package models for Book
package models

type Book struct {
	ID          int64
	Title       string
	Description string
	Authors     []string
	userID      int64
}
