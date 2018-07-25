package handlers

import (
	"github.com/AITestingOrg/calculation-service/controllers"
	"github.com/AITestingOrg/calculation-service/interfaces"
	"github.com/AITestingOrg/calculation-service/utils"
)

type MainProgram struct {
	interfaces.MainProgramInterface
}

func (m *MainProgram) BuildPublisher() *utils.AmqpPublisher {
	return new(utils.AmqpPublisher)
}

func (m *MainProgram) BuildConsumer(publisher *utils.AmqpPublisher) utils.AmqpConsumer {
	return utils.AmqpConsumer{
		"trip.exchange.tripcalculation",
		"topic",
		"trip.queue.calculationservice.calculatecost",
		"trip.estimation.estimatecalculated",
		EstimateHandler{publisher},
	}
}

func (m *MainProgram) BuildController(publisher *utils.AmqpPublisher) controllers.CalculationServiceController {
	apiHandlers := []interfaces.ApiHandlerInterface{
		CostEstimateHandler{Publisher: publisher}}
	return controllers.CalculationServiceController{Handlers: apiHandlers}
}

func (m *MainProgram) Run(publisher *utils.AmqpPublisher, controller controllers.CalculationServiceController, consumer utils.AmqpConsumer) {
	go controller.InitializeEndpoint()
	go utils.InitializeRabbitMqConsumers(consumer)
	go publisher.InitializePublisher()
}
