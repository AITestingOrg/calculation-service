package interfaces

import (
	"github.com/AITestingOrg/calculation-service/interfaces/data"
	"github.com/streadway/amqp"
)

type RabbitHandlerInterface interface {
	Handle(msg amqp.Delivery, session data.DataAccessInterface) error
}
