package rabbit

import (
	"github.com/Anacardo89/lenic/pkg/rabbitmq"
	"github.com/streadway/amqp"
)

func MQSendRegisterMail(r *rabbitmq.Config, ch *amqp.Channel, data []byte) error {
	err := ch.Publish(
		"",               // exchange
		"register_email", // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})
	if err != nil {
		return err
	}

	return nil
}

func MQSendPasswordRecoveryMail(r *rabbitmq.Config, ch *amqp.Channel, data []byte) error {
	err := ch.Publish(
		"",                       // exchange
		"password_recover_email", // routing key
		false,                    // mandatory
		false,                    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})
	if err != nil {
		return err
	}
	return nil
}
