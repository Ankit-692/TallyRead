// Package db for Database Creation
package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
)

var (
	DB          *sql.DB
	RedisClient *redis.Client
)

func InitDB() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	DB, err = sql.Open("pgx", connStr)
	if err != nil {
		panic("could not Open the db")
	}

	err = DB.Ping()
	if err != nil {
		panic("Connection is dead (Ping failed): " + err.Error())
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic("Could not parse Redis URL: " + err.Error())
	}

	RedisClient = redis.NewClient(opt)

	CreateTables()
}

func CreateTables() {
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        firstName TEXT NOT NULL,
        lastName TEXT,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        reset_token TEXT,
        reset_token_expires TIMESTAMP,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        search_count INTEGER DEFAULT 0,
        last_search_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	createBooksTable := `
    CREATE TABLE IF NOT EXISTS books (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        description TEXT,
        authors TEXT,
        user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        total_page INTEGER,
        ratings TEXT,
        image TEXT,
        published_date TEXT,
        page_read INTEGER CHECK (page_read >= 0),
        state TEXT,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(user_id, title)
    );`

	_, err := DB.Exec(createUsersTable)
	if err != nil {
		panic("could not create Users Table: " + err.Error())
	}

	_, err = DB.Exec(createBooksTable)
	if err != nil {
		panic("could not create Books Table: " + err.Error())
	}
}
