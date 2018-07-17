package utils

import (
	"github.com/streadway/amqp"
	"os"
	"log"
	"encoding/json"
)

var messagesChannel = make(chan Message)
var notifyCloseChannel =  make(chan bool)
var publishing = false

func PublishMessage(exchangeName string, routingKey string, payload interface{}) {
	//Todo add some validation here
	if(!publishing){
		//throw an error about how they need to start another go-routine and run StartPublishingMessagesFromChannel
	}
	encodedPayload, marshallErr := json.Marshal(payload)
	if marshallErr != nil {
		//do error handling
	}
	messagesChannel <- Message{exchangeName, routingKey, encodedPayload}
}

func StopPublishingMessagesFromChannel() {
	if publishing {
		notifyCloseChannel <- true
		publishing = false
	} else {
		//throw an error about how its not currently publishing
	}
}

func InitializeRabbitMqPublisher() {
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

	go func(){
		for message := range messagesChannel {
			err = ch.Publish(
				message.ExchangeName,
				message.RoutingKey,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:		 message.Message,
				})
			failOnError(err, "Failed to publish message")
			log.Printf("Published msg: " + string(message.Message) + " to exchange: " + message.ExchangeName + " with key: " + message.RoutingKey)
		}
	}()
	publishing = true
	<-notifyCloseChannel
}
