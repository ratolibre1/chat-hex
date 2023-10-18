package emitter

import (
	"github.com/streadway/amqp"
)

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) SendMessage(message string, ch *amqp.Channel, queueName string) error {
    err := ch.Publish(
        "",
        queueName,
        false,
        false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        },
    )
    return err
}