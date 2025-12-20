// Package db for Database Creation
package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "api.db")
	if err != nil {
		panic("could not Open the db")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	CreateTables()
}

func CreateTables() {
	createBooksTable := `
	CREATE TABLE IF NOT EXISTS books (
        id          INTEGER PRIMARY KEY,
        title       TEXT NOT NULL,
        description TEXT,
        authors     TEXT,
        user_id     INTEGER NOT NULL,
		total_page     TEXT,
        ratings        TEXT,
        image          TEXT,
        published_date TEXT,
        page_read      INTEGER CHECK (page_read IS NULL OR page_read >= 0),
		state       TEXT,
        created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    )		
	`

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)
	`

	_, err := DB.Exec(createBooksTable)
	if err != nil {
		panic("Could not create Books Table " + err.Error())
	}

	_, err = DB.Exec(createUsersTable)
	if err != nil {
		panic("could not create Users Table " + err.Error())
	}
}
