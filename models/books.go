// Package models for Book
package models

import (
	"time"

	"tallyRead.com/db"
)

type Book struct {
	ID            int64
	Title         string
	Description   string
	Authors       string
	TotalPage     string
	Ratings       string
	Image         string
	PublishedDate string
	PageRead      int16
	State         string
	UserID        int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (book *Book) Save() error {
	query := `
	INSERT INTO books(title,description,authors,user_id,total_page,ratings,image,published_date)
	VALUES (?,?,?,?,?,?,?,?)
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(book.Title, book.Description, book.Authors, book.UserID, book.TotalPage, book.Ratings, book.Image, book.PublishedDate)

	if err != nil {
		return err
	}

	book.ID, err = result.LastInsertId()

	return nil
}

func GetAllBooks(userId any) ([]Book, error) {
	query := `
	SELECT * FROM books WHERE user_id=?
	`
	rows, err := db.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books = []Book{}

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Description, &book.Authors, &book.UserID, &book.TotalPage, &book.Ratings, &book.Image, &book.PublishedDate, &book.PageRead, &book.State, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}
