package main

import (
	"github.com/AITestingOrg/calculation-service/eureka"
	"github.com/AITestingOrg/calculation-service/utils"
	"github.com/AITestingOrg/calculation-service/handlers"
	"github.com/AITestingOrg/calculation-service/controllers"
)

func main() {
	eureka.InitializeEurekaConnection()

	forever := make(chan bool)

	go controllers.InitializeEndpoint()

	go utils.InitializeRabbitMqConsumers(
		utils.Consumer{
		"trip.exchange.tripcalculation",
		"topic",
		"trip.queue.calculationservice.calculatecost",
		"trip.estimation.estimatecalculated",
		handlers.EstimateReceived,
		})

	go utils.InitializeRabbitMqPublisher()

	<- forever
}
