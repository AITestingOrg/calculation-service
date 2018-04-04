package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/AITestingOrg/calculation-service/models"
	"net/http"
	"log"
	"encoding/xml"
	"os"
	"time"
)

type Instance struct {
    IpAddress struct {
        InnerXML string `xml:",innerxml"`
    } `xml:"instance>ipAddr"`
}

func GetGmapsEstimation(trip models.Trip) models.Estimation {

	ipAddress := getIpAddress()

	// Encode trip
	encodedTrip, marshallErr := json.Marshal(trip)
	if marshallErr != nil {
		fmt.Println(marshallErr)
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
	var duration = float64(gmapsEstimation.Duration/60)
	var distance = float64(int(gmapsEstimation.Distance/1609 * 100)) / 100

	//Calculates cost
	log.Printf("Calculating costs...")
	var costDuration = (duration)*costPerMinute
	var costDistance = (distance)*costPerMile
	var finalCost = float64(int((costDuration + costDistance) * 100)) / 100

	//Maps response to JSON body
	currentDate := time.Now().Format("Jan 02 2006")
	encodedEstimation, marshallErr := json.Marshal(models.Estimation{ Cost: finalCost, 
		Duration: int64(duration), Distance: distance, Origin: trip.Origin, Destination: trip.Destination,
		LastUpdated: currentDate})
	if marshallErr != nil {
		fmt.Println(marshallErr)
		panic(marshallErr)
	}

	return encodedEstimation
}

func getIpAddress() string {

	eureka := os.Getenv("EUREKA_SERVER")
	if eureka == "" {
		eureka = "discovery-service"
	}
	url := fmt.Sprintf("http://%s:8761/eureka/apps/gmapsadapter", eureka)
	request, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	response, responseErr := client.Do(request)
	if responseErr != nil {
		panic(responseErr)
	}

	log.Printf("Reading XML body...")
	body, _ := ioutil.ReadAll(response.Body)

	var instance Instance
	unmarshallError := xml.Unmarshal(body, &instance)
	if unmarshallError != nil {
		panic(unmarshallError)
	}

	log.Printf("Received Ip Address!")

	return instance.IpAddress.InnerXML
}