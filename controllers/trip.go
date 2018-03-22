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
)

func GetCost(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	//Converts objects into json
	var trip models.Trip
	json.Unmarshal(body, &trip)
	var estimation models.Estimation
	json.Unmarshal(body, &estimation)

	// trip.ValidateOrigin(trip.origin)
	log.Printf("Validating trip and estimation body...")

	//Sets depature time to current time
	log.Printf("Sending current time to Gmaps adapter...")
	currentTime := time.Now().Unix()
	log.Print(currentTime)
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