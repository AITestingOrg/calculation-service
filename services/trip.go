package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AITestingOrg/calculation-service/models"
	"github.com/AITestingOrg/calculation-service/utils"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetGmapsEstimation(trip models.Trip) models.Estimation {

	ipAddress := utils.GetIpAddress()
	if ipAddress == "" {
		ipAddress = "localhost"
	}
	log.Printf("IP Address found: " + ipAddress)

	// Encode trip
	encodedTrip, marshallErr := json.Marshal(trip)
	if marshallErr != nil {
		panic(marshallErr)
	}

	// Request estimation
	log.Printf("Requesting information from Gmaps adapter...")
	url := fmt.Sprintf("http://%s:8080/api/v1/directions", ipAddress)

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(encodedTrip))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)
	if responseErr != nil {
		panic(responseErr)
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	// Decode response
	decodedEstimation := models.Estimation{}
	unmarshallError := json.Unmarshal(body, &decodedEstimation)
	if unmarshallError != nil {
		panic(unmarshallError)
	}
	log.Printf("Received information from Gmaps adapter!")
	return decodedEstimation
}

func CalculateCost(trip models.Trip, estimation models.Estimation) []byte {

	//Cost/Minute and Cost/Mile
	var costPerMinute = 0.15
	var costPerMile = 0.9

	//Get duration and distance from gmaps request
	gmapsEstimation := GetGmapsEstimation(trip)
	var duration = float64(gmapsEstimation.Duration / 60)
	var distance = float64(int(gmapsEstimation.Distance/1609*100)) / 100

	//Calculates cost
	log.Printf("Calculating costs...")
	var costDuration = (duration) * costPerMinute
	var costDistance = (distance) * costPerMile
	var finalCost = float64(int((costDuration+costDistance)*100)) / 100

	//Maps response to JSON body
	currentDate := time.Now().Format("2006-01-02 03:04:05")
	encodedEstimation, marshallErr := json.Marshal(models.Estimation{ Cost: finalCost, 
		Duration: int64(duration), Distance: distance, Origin: trip.Origin, Destination: trip.Destination,
		LastUpdated: currentDate})
	if marshallErr != nil {
		panic(marshallErr)
	}

	return encodedEstimation
}
