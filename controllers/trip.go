package controllers

import (
	"encoding/json"
	"errors"
	"github.com/AITestingOrg/calculation-service/models"
	"github.com/AITestingOrg/calculation-service/services"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
		err := errors.New("Strings are empty!")
		panic(err)
	}

	//Sets depature time to current time
	log.Printf("Sending current time to Gmaps adapter...")
	currentTime := time.Now().Unix()
	trip.DepartureTime = currentTime

	//Receives calculated cost and returns as json
	calculateCost := services.CalculateCost(trip, estimation)
	log.Printf("Distance, duration, and cost estimations returned!")

	w.Header().Set("Content-Type", "application/json")
	w.Write(calculateCost)
}
