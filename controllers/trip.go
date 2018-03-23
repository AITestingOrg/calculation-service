package controllers

import (
	"encoding/json"
	"io/ioutil"
	"calculation-service/services"
	"calculation-service/models"
	"net/http"
	"time"
	"fmt"
	"log"
	"errors"
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
	log.Printf("Receiving calculated costs...")
	encodedEstimationCost, marshallErr := json.Marshal(models.Cost{Cost: calculateCost})
	if marshallErr != nil {
		fmt.Println(marshallErr)
		panic(marshallErr)
	}
	log.Printf("Calculated costs returned!")
	w.Header().Set("Content-Type", "application/json")
	w.Write(encodedEstimationCost)
}