package utils

import (
	"log"
	"github.com/streadway/amqp"
	"os"
	"fmt"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)

		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func PublishMessage(ex_name string, binding string, msg []byte) {
	conn, err := amqp.Dial("amqp://guest:guest@" + os.Getenv("RABBIT_HOST") + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	log.Printf("after dialing rabbitmq")


	err = ch.Publish(
		ex_name,
		binding,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:		 msg,
		})
	failOnError(err, "Failed to publish message")
	log.Printf("Published msg: " + string(msg) + " to exchange: " + ex_name + " with key: " + binding)
}

func InitializeConsumer(ex_name string, ex_kind string, q_name string, q_binding string, handle func(msg amqp.Delivery)) {
	log.Printf("before dialing rabbitmq")

	conn, err := amqp.Dial("amqp://guest:guest@" + os.Getenv("RABBIT_HOST") + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	log.Printf("after dialing rabbitmq")

	err = ch.ExchangeDeclare(
		ex_name,
		ex_kind,
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare an exchange")

	log.Printf("after declaring exchange")

	q, err := ch.QueueDeclare(
		q_name,
		false,
		true,
		false,
		false,
		nil,
	)

	ch.QueueBind(q.Name,q_binding,ex_name,false, nil)

	failOnError(err, "Failed to declare the queue: " + q.Name )

	log.Printf("Declared queue: " + q.Name)

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
			handle(d)
		}
	}()

	log.Printf(" [*] Waiting for messages")
	<-forever
}