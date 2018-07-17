package utils

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

type Consumer struct {
	ExchangeName string
	ExchangeKind string
	QueueName string
	QueueBinding string
	Handle func(msg amqp.Delivery)
}

func (consumer Consumer) InitializeConsumer() {
	rabbitUsername := os.Getenv("RABBIT_USERNAME")
	if rabbitUsername == "" {
		rabbitUsername = "guest"
	}

	rabbitPassword := os.Getenv("RABBIT_PASSWORD")
	if rabbitPassword == "" {
		rabbitPassword = "guest"
	}

	rabbitHost := os.Getenv("RABBIT_HOST")
	if rabbitHost == "" {
		rabbitHost = "localhost"
	}

	conn, err := amqp.Dial("amqp://" + rabbitUsername + ":" + rabbitPassword + "@" + rabbitHost + ":5672/")
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
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			consumer.Handle(d)
		}
	}()

	log.Printf("Initialized queue (%s) bound to (%s) exchange, (%s) with binding key (%s)", consumer.QueueName, consumer.ExchangeKind, consumer.ExchangeName, consumer.QueueBinding)
	<-forever
}
