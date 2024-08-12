package rabbit

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitConfig struct {
	MQHost string `yaml:"mqHost"`
	MQPort string `yaml:"mqPort"`
}

var (
	RabbitMQ *RabbitConfig
)

func (r *RabbitConfig) MQSendRegisterMail(data []byte) error {
	rabbitUrl := fmt.Sprintf("amqp://%s%s", r.MQHost, r.MQPort)
	conn, err := amqp.Dial(rabbitUrl)
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := channel.QueueDeclare(
		"register_mail", // name
		true,            // durable
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

func (r *RabbitConfig) MQSendPasswordRecoveryMail(data []byte) error {
	rabbitUrl := fmt.Sprintf("amqp://%s%s", r.MQHost, r.MQPort)
	conn, err := amqp.Dial(rabbitUrl)
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := channel.QueueDeclare(
		"password_recover_mail", // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
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
