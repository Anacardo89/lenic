package mail

import (
	"fmt"
	"net/smtp"

	"github.com/Anacardo89/lenic/config"
)

const (
	mailMsg = `
	From: %s\r\n
	To: %s\r\n
	Subject: %s\r\n\r\n
	%s
	`
)

type Client struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewClient(cfg config.Mail) *Client {
	return &Client{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.User,
		Password: cfg.Pass,
	}
}

func (c *Client) Send(to []string, subject, body string) []error {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	auth := smtp.PlainAuth("", c.Username, c.Password, c.Host)
	var errs []error
	for _, t := range to {
		msg := []byte(fmt.Sprintf(mailMsg, c.Username, t, subject, body))
		err := smtp.SendMail(addr, auth, c.Username, []string{t}, msg)
		if err != nil {
			errs = append(errs, fmt.Errorf("error sending mail to: %s", t))
		}
	}
	return errs
}
