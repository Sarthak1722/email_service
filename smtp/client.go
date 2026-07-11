package smtp

import (
	"fmt"
	"os"

	mail "github.com/wneessen/go-mail"
)

type Client struct {
	client *mail.Client
	from   string
}

func New() (*Client, error) {

	host := os.Getenv("SMTP_HOST")
	port := 587
	user := os.Getenv("EMAIL_USERNAME")
	pass := os.Getenv("EMAIL_PASSWORD")

	c, err := mail.NewClient(
		host,
		mail.WithPort(port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(user),
		mail.WithPassword(pass),
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		client: c,
		from:   user,
	}, nil
}

func (c *Client) Send(to, subject, body string, isHTML bool) error {

	msg := mail.NewMsg()

	if err := msg.From(c.from); err != nil {
		return err
	}

	if err := msg.To(to); err != nil {
		return err
	}

	msg.Subject(subject)

	if isHTML {
		msg.SetBodyString(mail.TypeTextHTML, body)
	} else {
		msg.SetBodyString(mail.TypeTextPlain, body)
	}

	if err := c.client.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
