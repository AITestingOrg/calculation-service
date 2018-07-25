package handlers

import (
	"net/http"
	"io/ioutil"
	"github.com/AITestingOrg/calculation-service/models"
	"encoding/json"
	"log"
	"errors"
	"time"
	"github.com/AITestingOrg/calculation-service/interfaces"
)

type CostEstimateHandler struct {
	Publisher interfaces.PublisherInterface
}

func (controller CostEstimateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	//Converts objects into json
	var trip models.Trip
	json.Unmarshal(body, &trip)

	log.Printf("Validating trip and estimation body...")
	if (!trip.ValidateOrigin(trip.Origin)) || (!trip.ValidateDestination(trip.Destination)) {
		err := errors.New("ERROR: Invalid origin/destination address")
		log.Print(err)
		http.Error(w, err.Error(), 400)
		return
	}

	//Sets departure time to current time
	log.Printf("Setting current time for TripEstimate...")
	currentTime := time.Now().Unix()
	trip.DepartureTime = currentTime

	//Emit to trip.exchange.tripcalculation, with key trip.estimation.estimaterequested
	err := controller.Publisher.PublishMessage("trip.exchange.tripcalculation", "trip.estimation.estimaterequested", trip)
	if err == nil {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Trip Estimate Request Retrieved. Forwarding request to Gmaps Adapter"))
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(500)
		w.Write([]byte("Error publishing message to rabbitmq: " + err.Error()))
	}
}

func (handler CostEstimateHandler) GetPath() string{
	return "/api/v1/path"
}

func (handler CostEstimateHandler) GetRequestType() string{
	return "POST"
}
