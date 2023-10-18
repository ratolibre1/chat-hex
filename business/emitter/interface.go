package emitter

import "github.com/streadway/amqp"

type Service interface {
	SendMessage(message string, ch *amqp.Channel, queueName string) error
}