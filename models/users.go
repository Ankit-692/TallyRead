package models

import (
	"errors"
	"time"

	"tallyRead.com/db"
	"tallyRead.com/utils"
)

type User struct {
	ID                int64
	FirstName         string
	LastName          string
	Email             string `binding:"required"`
	Password          string `binding:"required"`
	ResetToken        string
	ResetTokenExpires time.Time
	CreatedAt         time.Time
}

func (u *User) Save() error {
	query := `
	INSERT INTO users (firstName,lastName,email,password)
	VALUES (?,?,?,?)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.FirstName, u.LastName, u.Email, hashedPassword)
	if err != nil {
		return err
	}

	u.ID, err = result.LastInsertId()
	u.Password = hashedPassword
	return err
}

func (u *User) ValidateUser() error {
	query := `
	SELECT id,email,firstName,lastName,password FROM users WHERE email=?
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
	SELECT id,email,firstName,lastName FROM users WHERE email=?
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
        SET reset_token = ?, reset_token_expires = ?
        WHERE id = ?
    `

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(u.ResetToken, u.ResetTokenExpires, u.ID)

	if err != nil {
		return err
	}

	return nil
}

func FindByResetToken(hashedToken string) (*User, error) {
	query := `
		SELECT id,email,firstName,lastName,reset_token,reset_token_expires 
		FROM users 
		WHERE reset_token=?
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
		SET password=?,reset_token=?,reset_token_expires=?
		WHERE email=?
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(hashedPassword, u.ResetToken, u.ResetTokenExpires, u.Email)

	if err != nil {
		return err
	}

	return nil

}
