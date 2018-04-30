package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/AITestingOrg/calculation-service/models"
	"github.com/AITestingOrg/calculation-service/services"
)

func GetCost(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	//Converts objects into json
	var trip models.Trip
	json.Unmarshal(body, &trip)
	var estimation models.Estimation
	json.Unmarshal(body, &estimation)

	log.Printf("Validating trip and estimation body...")
	if (!trip.ValidateOrigin(trip.Origin)) || (!trip.ValidateDestination(trip.Destination)) {
		err := errors.New("ERROR: Invalid origin/destination address")
		log.Print(err)
		http.Error(w, err.Error(), 400)
		return
	}

	//Sets depature time to current time
	log.Printf("Sending current time to Gmaps adapter...")
	currentTime := time.Now().Unix()
	trip.DepartureTime = currentTime

	//Receives calculated cost and returns as json
	calculateCost, err := services.CalculateCost(trip, estimation)
	if err != nil {
		err := errors.New("ERROR: Invalid trip request")
		log.Print(err)
		http.Error(w, err.Error(), 400)
		return
	}
	log.Printf("Distance, duration, and cost estimations returned!")

	w.Header().Set("Content-Type", "application/json")
	w.Write(calculateCost)
}
