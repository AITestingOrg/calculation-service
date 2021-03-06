package main

import (
	"github.com/AITestingOrg/calculation-service/eureka"
	"github.com/AITestingOrg/calculation-service/handlers"
	"github.com/AITestingOrg/calculation-service/interfaces"
	"github.com/AITestingOrg/calculation-service/utils"
)

func main() {
	//make an AmqpPublisher to be injected into the following methods
	amqpPublisher := new(utils.AmqpPublisher)

	//make a list of api handlers that should all be added to a http controller
	apiHandlers := []interfaces.ApiHandlerInterface{handlers.CostEstimateHandler{Publisher: amqpPublisher},
		handlers.GetCostHandler{}}

	//make a list of amqp consumers that should be consuming eventually
	amqpConsumers := []interfaces.ConsumerInterface{
		utils.AmqpConsumer{"trip.exchange.tripcalculation",
			"topic",
			"trip.queue.calculationservice.calculatecost",
			"trip.estimation.estimatecalculated",
			handlers.EstimateHandler{amqpPublisher},
		}}
	forever := make(chan bool)
	go utils.ProgramSetup(amqpPublisher, apiHandlers, amqpConsumers, eureka.EurekaClient{})
	<-forever
}
