package main

import (
	"github.com/AITestingOrg/calculation-service/eureka"
	"github.com/AITestingOrg/calculation-service/handlers"
)

func main() {
	eureka.InitializeEurekaConnection()

	main := new(handlers.MainProgram)
	forever := make(chan bool)
	amqpPublisher := main.BuildPublisher()
	controller := main.BuildController(amqpPublisher)
	consumer := main.BuildConsumer(amqpPublisher)
	main.Run(amqpPublisher, controller, consumer)

	<-forever
}
