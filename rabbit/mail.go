package rabbit

import (
	"github.com/streadway/amqp"
)

type RabbitConfig struct {
	MQHost string `yaml:"mqHost"`
	MQPort string `yaml:"mqPort"`
}

var (
	RabbitMQ *RabbitConfig
)

func (r *RabbitConfig) MQSendRegMail(data []byte) error {
	url := "amqp://" + r.MQHost + r.MQPort
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := channel.QueueDeclare(
		"register_mail", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return err
	}

	err = channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})
	if err != nil {
		return err
	}

	return nil
}
