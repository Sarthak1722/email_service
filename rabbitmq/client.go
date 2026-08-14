package rabbitmq

import (
	"encoding/json"
	"os"

	"github.com/Sarthak1722/email_service/email"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}
const QueueName = "email_jobs"

func New() (*Client, error) {

	url := os.Getenv("RABBITMQ_URL")

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = channel.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}
	return &Client{
		conn:    conn,
		channel: channel,
	}, nil
}

func (c *Client) Publish(job email.Job) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return c.channel.Publish(
    "",
    QueueName,
    false,
    false,
    amqp.Publishing{
        ContentType: "application/json",
        Body:        body,
    },
)
}

func (c *Client) Consume() (<-chan amqp.Delivery, error) {

    return c.channel.Consume(
        QueueName,
        "",
        false,
        false,
        false,
        false,
        nil,
    )
}