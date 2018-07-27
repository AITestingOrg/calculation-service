package interfaces

import "github.com/streadway/amqp"

type RabbitHandlerInterface interface {
	Handle(msg amqp.Delivery) error
}
