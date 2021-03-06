package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/AITestingOrg/calculation-service/interfaces"
	"github.com/AITestingOrg/calculation-service/interfaces/data"
	"github.com/AITestingOrg/calculation-service/models"
	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2/bson"
)

type EstimateHandler struct {
	Publisher interfaces.PublisherInterface
}

func (handler EstimateHandler) Handle(msg amqp.Delivery, session data.DataAccessInterface) error {
	genericMessageReceived(msg)
	data := msg.Body
	var estimation models.Estimation
	err := json.Unmarshal(data, &estimation)
	if err != nil {
		return errors.New("error unmarshalling data into an estimation object: " + err.Error())
	}

	err = estimation.ValidateFields("origin", "destination", "distance", "duration", "userId")
	if err != nil {
		return errors.New("error with the parsed estimation object: \n" + err.Error())
	}

	estimate := calculateCost(estimation, session)

	err = handler.Publisher.PublishMessage("notification.exchange.notification", "notification.trip.estimatecalculated", estimate)
	if err != nil {
		return errors.New("error publishing message: " + err.Error())
	}

	return nil
}

func genericMessageReceived(msg amqp.Delivery) {
	log.Printf("Message received on exchange: %s\n\tWith routing key: %s\n\tWith body: %s", msg.Exchange, msg.RoutingKey, msg.Body)
}

func calculateCost(gmapsEstimation models.Estimation, session data.DataAccessInterface) models.Estimation {
	
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

	//Create cost object to be saved
	log.Println("Writing cost to database")
	c := session.C("costs")

	if c != nil {
		_, err2 := c.Upsert(bson.M{"userId": gmapsEstimation.UserId}, models.Cost{
			UserId:      gmapsEstimation.UserId,
			Origin:      gmapsEstimation.Origin,
			Destination: gmapsEstimation.Destination,
			Cost:        finalCost,
		})

		if err2 != nil {
			log.Printf("Error saving to database: %s", err2)
		}
	} else {
		log.Printf("No registries found in database.")
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
