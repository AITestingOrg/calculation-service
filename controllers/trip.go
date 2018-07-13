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

//need to change this to take in trip and user id
func GetCost(w http.ResponseWriter, r *http.Request) {
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

	var tripEstimate models.TripEstimateRequest
	json.Unmarshal(body, &tripEstimate)

	//Sets depature time to current time
	log.Printf("Setting current time for TripEstimate...")
	currentTime := time.Now().Unix()
	tripEstimate.DepartureTime = currentTime

	services.EmitGmapsEstimationRequest(tripEstimate)

	//
	////Receives calculated cost and returns as json
	//if err != nil {
	//	err := errors.New("ERROR: Invalid trip request")
	//	log.Print(err)
	//	http.Error(w, err.Error(), 400)
	//	return
	//}
	//log.Printf("Distance, duration, and cost estimations returned!")
	//
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Trip Estimate Request Retrieved. Forwarding request to Gmaps Adapter"))
}
