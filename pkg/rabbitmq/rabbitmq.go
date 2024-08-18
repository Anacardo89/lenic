package rabbitmq

import (
	"fmt"

	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
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

func (r *Config) Connect() *amqp.Connection {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		r.RabbitUser, r.RabbitPass, r.RabbitHost, r.RabbitPort)
	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Error.Fatal(err)
	}
	return conn
}

func (r *Config) DeclareQueues(ch *amqp.Channel) {
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
			logger.Error.Fatal(err)
		}
	}
}
