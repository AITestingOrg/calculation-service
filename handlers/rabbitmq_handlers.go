package handlers

import (
	"github.com/streadway/amqp"
	"github.com/AITestingOrg/calculation-service/models"
	"encoding/json"
	"log"
	"github.com/AITestingOrg/calculation-service/utils"
	"time"
)

func GenericMessageReceived(msg amqp.Delivery){
	log.Printf("Message received on exchange: %s\n\tWith routing key: %s\n\tWith body: %s", msg.Exchange, msg.RoutingKey, msg.Body)
}

func EstimateReceived(msg amqp.Delivery){
	GenericMessageReceived(msg)
	data := msg.Body
	var estimation models.Estimation
	json.Unmarshal(data, &estimation)
	estimate := CalculateCost(estimation)
	log.Printf("Message received and unmarshaled into an estimation object: %s", estimation)
	utils.PublishMessage("notification.exchange.notification", "notification.trip.estimatecalculated", estimate)
	msg.Ack(false)
}

func CalculateCost(gmapsEstimation models.Estimation) (models.Estimation) {
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
