package handlers

import (
	"encoding/json"
	"errors"
	"github.com/AITestingOrg/calculation-service/interfaces"
	"github.com/AITestingOrg/calculation-service/models"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
)

type CostEstimateHandler struct {
	Publisher interfaces.PublisherInterface
}

func (handler CostEstimateHandler) AddHandlerToRouter(r * mux.Router) {
	r.HandleFunc("/api/v1/cost", handler.Handle).Methods("POST")
}

func (handler CostEstimateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	//Converts objects into json
	var trip models.Trip
	json.Unmarshal(body, &trip)

	log.Printf("Validating trip and estimation body...")
	err := trip.ValidateFields("originAddress", "destinationAddress", "userId")
	if err != nil {
		err = errors.New("ERROR: Invalid trip arguments:\n" + err.Error())
		log.Print(err)
		http.Error(w, err.Error(), 400)
		return
	}

	//Sets departure time to current time
	log.Printf("Setting current time for TripEstimate...")
	currentTime := time.Now().Unix()
	trip.DepartureTime = currentTime

	//Emit to trip.exchange.tripcalculation, with key trip.estimation.estimaterequested
	err = handler.Publisher.PublishMessage("trip.exchange.tripcalculation", "trip.estimation.estimaterequested", trip)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(500)
		w.Write([]byte("Error publishing message to rabbitmq: " + err.Error()))
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Trip Estimate Request Retrieved. Forwarding request to Gmaps Adapter"))
}
