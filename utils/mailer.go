package utils

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/mail.v2"
	"tallyRead.com/config"
	"tallyRead.com/templates"
)

func SendResetEmail(to, resetURL string) error {
	configSmtp := config.GetSMTP()
	port, _ := strconv.Atoi(configSmtp.Port)
	m := mail.NewMessage()
	m.SetHeader("From", configSmtp.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Reset your password")

	htmlBody := fmt.Sprintf(templates.PasswordResetTemplate, resetURL)
	m.SetBody("text/html", htmlBody)

	d := mail.NewDialer(configSmtp.Host, port, configSmtp.Email, configSmtp.Password)

	d.Timeout = 30 * time.Second
	d.TLSConfig = &tls.Config{ServerName: configSmtp.Host}

	err := d.DialAndSend(m)

	if err != nil {
		return err
	}

	return nil
}
