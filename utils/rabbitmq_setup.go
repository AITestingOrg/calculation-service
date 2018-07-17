package utils

import (
	"log"
	"fmt"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)

		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func InitializeRabbitMqConsumers(consumers ...Consumer){
	for _, consumer := range consumers{
		go consumer.InitializeConsumer()
	}
}
