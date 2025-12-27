package utils

import (
	"fmt"
	"strconv"

	"gopkg.in/mail.v2"
	"tallyRead.com/config"
	"tallyRead.com/templates"
)

func SendResetEmail(to, resetURL string) error {
	port, _ := strconv.Atoi(config.Smtp.Port)
	m := mail.NewMessage()
	m.SetHeader("From", config.Smtp.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Reset your password")

	htmlBody := fmt.Sprintf(templates.PasswordResetTemplate, resetURL)
	m.SetBody(htmlBody, resetURL)

	d := mail.NewDialer(config.Smtp.Host, port, config.Smtp.Email, config.Smtp.Password)

	err := d.DialAndSend(m)

	if err != nil {
		return err
	}

	return nil
}
