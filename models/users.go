package models

import (
	"errors"

	"tallyRead.com/db"
	"tallyRead.com/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := `
	INSERT INTO users (email,password)
	VALUES (?,?)
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

	result, err := stmt.Exec(u.Email, hashedPassword)
	if err != nil {
		return err
	}

	u.ID, err = result.LastInsertId()
	u.Password = hashedPassword
	return err
}

func (u *User) ValidateUser() error {
	query := `
	SELECT id,password FROM users WHERE email=?
	`

	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string

	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		return nil
	}

	isValid := utils.CheckPassword(u.Password, retrievedPassword)

	if !isValid {
		return errors.New("invalid credentials")
	}

	return nil
}
