// Package models for Book
package models

type Book struct {
	ID          int
	Title       string
	Description string
	Authors     []string
	userID      int64
}
