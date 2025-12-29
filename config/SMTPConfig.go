package config

import "os"

type SMTPConfig struct {
	Host      string
	Port      string
	Email     string
	Password  string
	ClientURL string
}

var smtpConfig *SMTPConfig

func GetSMTP() SMTPConfig {
	if smtpConfig == nil {
		smtpConfig = &SMTPConfig{
			Host:      os.Getenv("SMTP_HOST"),
			Port:      os.Getenv("SMTP_PORT"),
			Email:     os.Getenv("SMTP_EMAIL"),
			Password:  os.Getenv("SMTP_PASSWORD"),
			ClientURL: os.Getenv("FRONTEND_URL"),
		}
	}
	return *smtpConfig
}
