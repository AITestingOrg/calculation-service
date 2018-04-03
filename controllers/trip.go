package controllers

import (
	"fmt"
	"log"
	"time"
	"errors"
	"net/http"
	"io/ioutil"
	"encoding/json"
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
		err := errors.New("Strings are empty!")
		panic(err)
	}

	//Sets depature time to current time
	log.Printf("Sending current time to Gmaps adapter...")
	currentTime := time.Now().Unix()
	trip.DepartureTime = currentTime

	//Receives calculated cost and returns as json
	calculateCost := services.CalculateCost(trip, estimation)
	gmapsEstimation := services.GetGmapsEstimation(trip)
	duration := gmapsEstimation.Duration/60
	distance := float64(int(gmapsEstimation.Distance/1609 * 100)) / 100
	log.Printf("Receiving calculated costs...")
	encodedEstimation, marshallErr := json.Marshal(models.Estimation{ Cost: calculateCost, 
		Duration: duration, Distance: distance})
	if marshallErr != nil {
		fmt.Println(marshallErr)
		panic(marshallErr)
	}
	log.Printf("Distance, duration, and cost estimations returned!")
	w.Header().Set("Content-Type", "application/json")
	w.Write(encodedEstimation)
}