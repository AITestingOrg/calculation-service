package main

import (
	"github.com/AITestingOrg/calculation-service/controllers"
	"github.com/AITestingOrg/calculation-service/eureka"
	"github.com/AITestingOrg/calculation-service/handlers"
	"github.com/AITestingOrg/calculation-service/utils"
	"github.com/AITestingOrg/calculation-service/interfaces"
)

func main() {
	eureka.InitializeEurekaConnection()

	amqpPublisher := new(utils.AmqpPublisher)

	forever := make(chan bool)
	apiHandlers := []interfaces.ApiHandlerInterface{
		controllers.CostEstimateHandler{amqpPublisher, "/api/v1/cost", "POST"}}

	controller := controllers.CalculationServiceController{apiHandlers}
	go controller.InitializeEndpoint()

	go utils.InitializeRabbitMqConsumers(
		utils.AmqpConsumer{
		"trip.exchange.tripcalculation",
		"topic",
		"trip.queue.calculationservice.calculatecost",
		"trip.estimation.estimatecalculated",
		handlers.EstimateHandler{amqpPublisher},
		})

	go new(utils.AmqpPublisher).InitializePublisher()

	<- forever
}
