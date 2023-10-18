package listener

import (
	"github.com/streadway/amqp"
)

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) ConsumeMessages(ch *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
    msgs, err := ch.Consume(
        queueName,
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    return msgs, err
}