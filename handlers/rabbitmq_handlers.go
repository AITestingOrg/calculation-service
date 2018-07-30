package handlers

import (
	"encoding/json"
	"errors"
	"github.com/AITestingOrg/calculation-service/interfaces"
	"github.com/AITestingOrg/calculation-service/models"
	"github.com/streadway/amqp"
	"log"
	"time"
	"github.com/AITestingOrg/calculation-service/db"
	"gopkg.in/mgo.v2"
)

type EstimateHandler struct {
	Publisher interfaces.PublisherInterface
}

func (handler EstimateHandler) Handle(msg amqp.Delivery) error {
	genericMessageReceived(msg)
	data := msg.Body
	var estimation models.Estimation
	err := json.Unmarshal(data, &estimation)
	if err != nil {
		return errors.New("error unmarshalling data into an estimation object: " + err.Error())
	}

	err = estimation.ValidateFields("originAddress", "destinationAddress", "distance", "duration", "userId")
	if err != nil {
		return errors.New("error with the parsed estimation object: \n" + err.Error())
	}

	estimate := calculateCost(estimation)

	err = handler.Publisher.PublishMessage("notification.exchange.notification", "notification.trip.estimatecalculated", estimate)
	if err != nil {
		return errors.New("error publishing message: " + err.Error())
	}

	return nil
}

func genericMessageReceived(msg amqp.Delivery) {
	log.Printf("Message received on exchange: %s\n\tWith routing key: %s\n\tWith body: %s", msg.Exchange, msg.RoutingKey, msg.Body)
}

func calculateCost(gmapsEstimation models.Estimation) models.Estimation {
	//Copy Mongo session
	session := db.MgoSession.Copy()
	defer session.Close()
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

	//Set cost
	var cost models.Cost
	cost.Origin = gmapsEstimation.Origin
	cost.UserId = gmapsEstimation.UserId
	cost.DepartureTime = time.Now().Unix()
	cost.Destination = gmapsEstimation.Destination
	cost.Cost = finalCost
		//End set cost
		log.Println("Writing cost to database")
		c := session.DB("TRIPCOST").C("costs")
		log.Printf(cost.Destination, cost.DepartureTime)
		err := c.Insert(cost)
		if err != nil {
			if mgo.IsDup(err) {
				log.Printf("Error saving to database: %s", err)
			}
		}

	return models.Estimation{
		Cost:        finalCost,
		Duration:    int64(duration),
		Distance:    distance,
		Origin:      gmapsEstimation.Origin,
		Destination: gmapsEstimation.Destination,
		LastUpdated: currentDate,
		UserId:      gmapsEstimation.UserId,
	}
}

