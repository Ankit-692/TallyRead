package models

import (
	"errors"
	"fmt"
	"time"

	"tallyRead.com/db"
	"tallyRead.com/utils"
)

type User struct {
	ID                int64
	FirstName         string
	LastName          string
	Email             string `binding:"required"`
	Password          string `binding:"required,min=8"`
	ResetToken        string
	ResetTokenExpires time.Time
	CreatedAt         time.Time
	SearchCount       int64
	LastSearchDate    time.Time
}

func (u *User) Save() error {
	var count int
	countQuery := `SELECT COUNT(*) FROM users`
	err := db.DB.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return err
	}

	if count >= 350 {
		return fmt.Errorf("registration limit reached: maximum 350 users allowed")
	}

	query := `
		INSERT INTO users (firstName, lastName, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	err = db.DB.QueryRow(
		query,
		u.FirstName,
		u.LastName,
		u.Email,
		hashedPassword,
	).Scan(&u.ID)

	if err != nil {
		return err
	}

	u.Password = hashedPassword
	return nil
}

func (u *User) ValidateUser() error {
	query := `
	SELECT id,email,firstName,lastName,password FROM users WHERE email=$1
	`

	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string

	err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &retrievedPassword)
	if err != nil {
		return err
	}

	isValid := utils.CheckPassword(u.Password, retrievedPassword)

	u.Password = ""

	if !isValid {
		return errors.New("invalid credentials")
	}

	return nil
}

func FindByEmail(email string) (*User, error) {
	query := `
	SELECT id,email,firstName,lastName FROM users WHERE email=$1
	`
	row := db.DB.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) UpdateResetToken() error {
	query := `
        UPDATE users 
        SET reset_token = $1, reset_token_expires = $2
        WHERE id = $3
    `

	_, err := db.DB.Exec(query, u.ResetToken, u.ResetTokenExpires, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func FindByResetToken(hashedToken string) (*User, error) {
	query := `
		SELECT id,email,firstName,lastName,reset_token,reset_token_expires 
		FROM users 
		WHERE reset_token=$1
	`
	row := db.DB.QueryRow(query, hashedToken)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.ResetToken, &user.ResetTokenExpires)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) UpdatePassword() error {
	query := `
		UPDATE users
		SET password=$1,reset_token=$2,reset_token_expires=$3
		WHERE email=$4
	`
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	_, err = db.DB.Exec(query, hashedPassword, u.ResetToken, u.ResetTokenExpires, u.Email)
	if err != nil {
		return err
	}

	return nil
}

func FindSearchCountByID(userId int64) (*User, error) {

	query := `
		SELECT id,search_count, last_search_date 
		FROM users 
		WHERE id = $1
	`

	row := db.DB.QueryRow(query, userId)

	var user User

	err := row.Scan(&user.ID, &user.SearchCount, &user.LastSearchDate)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (u *User) UpdateSearchCount() error {

	query := `
		UPDATE users 
		SET search_count = $1, last_search_date = $2 
		WHERE id = $3	
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(u.SearchCount, u.LastSearchDate, u.ID)

	if err != nil {
		return err
	}

	return nil
}
