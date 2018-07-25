package utils

import (
	"github.com/streadway/amqp"
	"log"
	"github.com/AITestingOrg/calculation-service/interfaces"
)

type AmqpConsumer struct {
	ExchangeName string
	ExchangeKind string
	QueueName string
	QueueBinding string
	Handler interfaces.RabbitHandlerInterface
}

func (consumer AmqpConsumer) InitializeConsumer() {
	rabbitCreds := GetRabbitCredentials()

	conn, err := amqp.Dial("amqp://" + rabbitCreds["username"] + ":" + rabbitCreds["password"] + "@" + rabbitCreds["host"] + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		consumer.ExchangeName,
		consumer.ExchangeKind,
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		consumer.QueueName,
		false,
		true,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare the queue: " + q.Name )

	ch.QueueBind(consumer.QueueName,consumer.QueueBinding,consumer.ExchangeName,false, nil)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			err := consumer.Handler.Handle(msg)
			if err != nil {
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}
	}()

	log.Printf("Initialized queue: (%s) bound to %s-type exchange: (%s) with binding key: (%s)", consumer.QueueName, consumer.ExchangeKind, consumer.ExchangeName, consumer.QueueBinding)
	<-forever
}
