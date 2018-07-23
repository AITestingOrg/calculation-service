package handlers

import (
	"github.com/streadway/amqp"
	"github.com/AITestingOrg/calculation-service/models"
	"encoding/json"
	"log"
	"time"
	"errors"
	"github.com/AITestingOrg/calculation-service/interfaces"
)

type RabbitHandlerInterface interface {
	Handle(msg amqp.Delivery) error
}

type EstimateHandler struct {
	Publisher interfaces.PublisherInterface
}

func (handler EstimateHandler) Handle(msg amqp.Delivery) error{
	genericMessageReceived(msg)
	data := msg.Body
	var estimation models.Estimation
	json.Unmarshal(data, &estimation)
	estimate := calculateCost(estimation)
	err := handler.Publisher.PublishMessage("notification.exchange.notification", "notification.trip.estimatecalculated", estimate)
	if err == nil {
		return nil
	} else {
		return errors.New("error publishing message: " + err.Error())
	}
}

func genericMessageReceived(msg amqp.Delivery){
	log.Printf("Message received on exchange: %s\n\tWith routing key: %s\n\tWith body: %s", msg.Exchange, msg.RoutingKey, msg.Body)
}

func calculateCost(gmapsEstimation models.Estimation) (models.Estimation) {
	//Cost/Minute and Cost/Mile
	var costPerMinute = 0.15
	var costPerMile = 0.9

	//Get duration and distance from gmaps request
	var duration = float64(gmapsEstimation.Duration / 60)
	var distance = float64(int(gmapsEstimation.Distance/1609*100)) / 100

	//Calculates cost
	log.Printf("Calculating costs...")
	var costDuration = (duration) * costPerMinute
	var costDistance = (distance) * costPerMile
	var finalCost = float64(int((costDuration+costDistance)*100)) / 100

	currentDate := time.Now().Format("2006-01-02 03:04:05")

	return models.Estimation{
		Cost: finalCost,
		Duration: int64(duration),
		Distance: distance,
		Origin: gmapsEstimation.Origin,
		Destination: gmapsEstimation.Destination,
		LastUpdated: currentDate,
		UserId: gmapsEstimation.UserId,
	}
}
