// Package models for Book
package models

import (
	"time"

	"tallyRead.com/db"
)

type Book struct {
	ID            int64
	Title         string `json:"title"`
	Description   string `json:"description"`
	Authors       string `json:"authors"`
	TotalPage     int16  `json:"totalPage"`
	Ratings       string `json:"ratings"`
	Image         string `json:"image"`
	PublishedDate string `json:"publishedDate"`
	PageRead      int16
	State         string
	UserID        int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (book *Book) Save() error {
	query := `
	INSERT INTO books(title,description,authors,user_id,total_page,ratings,image,published_date,page_read,state)
	VALUES (?,?,?,?,?,?,?,?,?,?)
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(book.Title, book.Description, book.Authors, book.UserID, book.TotalPage, book.Ratings, book.Image, book.PublishedDate, book.PageRead, book.State)

	if err != nil {
		return err
	}

	book.ID, err = result.LastInsertId()

	return nil
}

func GetAllBooks(userId any) ([]Book, error) {
	query := `
	SELECT id, title, description, authors, total_page, ratings, image, 
	       published_date, page_read, state, user_id, created_at, updated_at 
	FROM books 
	WHERE user_id = ?
	`
	rows, err := db.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books = []Book{}

	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Description,
			&book.Authors,
			&book.TotalPage,
			&book.Ratings,
			&book.Image,
			&book.PublishedDate,
			&book.PageRead,
			&book.State,
			&book.UserID,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func GetBookByID(id int64) (*Book, error) {
	query := `SELECT * from books WHERE id = ?`

	row := db.DB.QueryRow(query, id)

	var book Book
	err := row.Scan(&book.ID, &book.Title, &book.Description, &book.Authors, &book.UserID, &book.TotalPage, &book.Ratings, &book.Image, &book.PublishedDate, &book.PageRead, &book.State, &book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &book, err
}

func (book *Book) Delete() error {
	query := `
	DELETE FROM books WHERE ID = ?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(book.ID)

	return err
}

func (book *Book) Update() error {
	query := `
	UPDATE books
	SET page_read=?, state=?, updated_at=CURRENT_TIMESTAMP 
	WHERE id=?
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(book.PageRead, book.State, book.ID)

	if err != nil {
		return err
	}

	return nil
}

func GetBooksByStatus(userId int64, status string) ([]Book, error) {
	query := `
	SELECT * FROM books
	WHERE user_id=? AND state=?
	`

	rows, err := db.DB.Query(query, userId, status)

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
