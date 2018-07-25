package utils

import (
	"fmt"
	"github.com/AITestingOrg/calculation-service/interfaces"
	"log"
	"os"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)

		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func InitializeRabbitMqConsumers(consumers ...interfaces.ConsumerInterface) {
	for _, consumer := range consumers {
		go consumer.InitializeConsumer()
	}
}

func GetRabbitCredentials() map[string]string {
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
	credMap := make(map[string]string)
	credMap["username"] = rabbitUsername
	credMap["password"] = rabbitPassword
	credMap["host"] = rabbitHost
	return credMap
}
