package listener

import "github.com/streadway/amqp"

type Service interface {
	ConsumeMessages(ch *amqp.Channel, queueName string) (<-chan amqp.Delivery, error)
}