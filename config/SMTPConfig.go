package config

import "os"

type SMTPConfig struct {
	Host      string
	Port      string
	Email     string
	Password  string
	ClientURL string
}

var Smtp = SMTPConfig{
	Host:      os.Getenv("SMTP_HOST"),
	Port:      os.Getenv("SMTP_PORT"),
	Email:     os.Getenv("SMTP_EMAIL"),
	Password:  os.Getenv("your-gmail-app-password"),
	ClientURL: os.Getenv("FRONTEND_URL"),
}
