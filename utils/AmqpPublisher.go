package utils

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"log"
)

type AmqpPublisher struct {
	publishingChannel *amqp.Channel
	closeSignal       chan bool
}

var (
	staticPublisher = &AmqpPublisher{publishingChannel: nil, closeSignal: make(chan bool)}
)

func (publisher *AmqpPublisher) PublishMessage(exchangeName string, routingKey string, payload interface{}) error {
	//Todo add some validation here
	if staticPublisher.publishingChannel == nil {
		return errors.New("publisher is not initialized, please call InitializePublisher() first")
	}
	encodedPayload, marshallErr := json.Marshal(payload)
	if marshallErr != nil {
		return errors.New("message was not published due to error marshalling payload: " + marshallErr.Error())
	}
	err := staticPublisher.publishingChannel.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        encodedPayload,
		})
	if err != nil {
		return errors.New("Error publishing message: " + err.Error())
	}
	log.Printf("Published msg: " + string(encodedPayload) + " to exchange: " + exchangeName + " with key: " + routingKey)
	return nil
}

func (publisher *AmqpPublisher) StopPublisher() error {
	if staticPublisher.publishingChannel != nil {
		log.Printf("Sending signal to close channel and stop publisher")
		staticPublisher.closeSignal <- true
		return nil
	} else {
		return errors.New("error stopping publisher, publisher not currently running")
	}
}

func (publisher *AmqpPublisher) InitializePublisher() error {
	if staticPublisher.publishingChannel == nil {
		log.Printf("Initializing publisher")

		rabbitCreds := GetRabbitCredentials()

		conn, err := amqp.Dial("amqp://" + rabbitCreds["username"] + ":" + rabbitCreds["password"] + "@" + rabbitCreds["host"] + ":5672/")
		if err != nil {
			return errors.New("failed to connect to rabbitmq: " + err.Error())
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			return errors.New("failed to open channel on rabbitmq connection: " + err.Error())
		}
		defer ch.Close()

		staticPublisher.publishingChannel = ch

		<-staticPublisher.closeSignal //blocking until StopPublisher is called

		staticPublisher.publishingChannel = nil
		return nil
	} else {
		return errors.New("error initializing publisher, publisher already initialized")
	}
}
