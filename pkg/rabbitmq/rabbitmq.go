package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Config struct {
	RabbitUser string   `yaml:"rabbit_user"`
	RabbitPass string   `yaml:"rabbit_pass"`
	RabbitHost string   `yaml:"rabbit_host"`
	RabbitPort string   `yaml:"rabbit_port"`
	Queues     []string `yaml:"queues"`
}

var (
	RMQ *Config
	RCh *amqp.Channel
)

func (r *Config) Connect() (*amqp.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		r.RabbitUser, r.RabbitPass, r.RabbitHost, r.RabbitPort)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (r *Config) DeclareQueues(ch *amqp.Channel) error {
	for _, queue := range r.Queues {
		_, err := ch.QueueDeclare(
			queue, // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return err
		}
	}
	return nil
}
