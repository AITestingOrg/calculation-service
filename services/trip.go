package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/AITestingOrg/calculation-service/models"
	"net/http"
	"log"
)

func GetGmapsEstimation(trip models.Trip) models.Estimation {
	// Encode trip
	encodedTrip, marshallErr := json.Marshal(trip)
	if marshallErr != nil {
		fmt.Println(marshallErr)
		panic(marshallErr)
	}

	// Request estimation
	log.Printf("Requesting information from Gmaps adapter...")
	url := "http://localhost:8080/api/v1/directions"

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

func CalculateCost(trip models.Trip, estimation models.Estimation) float64 {
	//Cost/Minute and Cost/Mile
	var costPerMinute = 0.15
	var costPerMile = 0.9

	//Get duration and distance from gmaps request
	gmapEstimation := GetGmapsEstimation(trip)
	var duration = float64(gmapEstimation.Duration)
	var distance = float64(gmapEstimation.Distance)

	//Calculates cost
	log.Printf("Calculating costs...")
	var costDuration = (duration/60)*costPerMinute
	var costDistance = (distance/1609.34)*costPerMile

	return float64(int((costDuration + costDistance) * 100)) / 100
}