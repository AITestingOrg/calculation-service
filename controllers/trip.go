package controllers

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"calculation-service/models"
	"calculation-service/services"
)

func GetCost(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	//Converts objects into json
	var trip models.Trip
	json.Unmarshal(body, &trip)
	var estimation models.Estimation
	json.Unmarshal(body, &estimation)
	
	log.Printf("Validating trip and estimation body...")

	//Sets depature time to current time
	log.Printf("Sending current time to Gmaps adapter...")
	currentTime := time.Now().Unix()
	log.Print(currentTime)
	trip.DepartureTime = currentTime

	//Receives calculated cost and returns as json
	calculateCost := services.CalculateCost(trip, estimation)
	log.Printf("Receiving calculated costs...")
	encodedEstimationCost, marshallErr := json.Marshal(models.Estimation{ Cost: calculateCost })
	if marshallErr != nil {
		fmt.Println(marshallErr)
		panic(marshallErr)
	}
	log.Printf("Distance, duration, and cost estimations returned!")
	w.Header().Set("Content-Type", "application/json")
	w.Write(encodedEstimationCost)
}