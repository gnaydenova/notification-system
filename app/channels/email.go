package channels

import (
	"fmt"
	"net/smtp"
	"strings"
)

type EmailConfig struct {
	Enabled  bool     `yaml:"enabled"`
	Host     string   `yaml:"host"`
	Port     string   `yaml:"port"`
	User     string   `yaml:"user"`
	Password string   `yaml:"password"`
	From     string   `yaml:"from"`
	To       []string `yaml:"to"`
	Subject  string   `yaml:"subject"`
}

type Email struct {
	config EmailConfig
}

func NewEmail(cfg EmailConfig) Email {
	return Email{config: cfg}
}

func (e Email) Send(msg string) error {
	if !e.config.Enabled {
		return nil
	}

	message := fmt.Sprintf("From: %s\r\n", e.config.From)
	message += fmt.Sprintf("To: %s\r\n", strings.Join(e.config.To, ","))
	message += fmt.Sprintf("Subject: %s\r\n", e.config.Subject)
	message += fmt.Sprintf("\r\n%s\r\n", msg)

	auth := smtp.PlainAuth("", e.config.User, e.config.Password, e.config.Host)

	err := smtp.SendMail(
		e.config.Host+":"+e.config.Port,
		auth,
		e.config.From,
		e.config.To,
		[]byte(message),
	)

	return err
}
