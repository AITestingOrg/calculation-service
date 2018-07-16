package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AITestingOrg/calculation-service/controllers"
	"github.com/AITestingOrg/calculation-service/eureka"
	"github.com/AITestingOrg/calculation-service/utils"

	"github.com/gorilla/mux"

	"fmt"
	"github.com/streadway/amqp"
	"github.com/AITestingOrg/calculation-service/models"
	"encoding/json"
	"github.com/AITestingOrg/calculation-service/services"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)

		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/cost", controllers.GetCost).Methods("POST")
	log.Println("Calculation service is running...")

	//Check to see if running locally or not
	var localRun = false
	if os.Getenv("EUREKA_SERVER") == "" {
		localRun = true
	}
	if !localRun {
		var eurekaUp = false
		log.Println("Waiting for Eureka...")
		for eurekaUp != true {
			eurekaUp = checkEurekaService(eurekaUp)
		}
		eureka.PostToEureka()
		eureka.StartHeartbeat()
		log.Printf("After scheduling heartbeat")
	}

	forever := make(chan bool)

	http.Handle("/", r)
	log.Printf("After http.handle")

	go func() {
		log.Printf("inside:: before listen and serve")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	genericMessageReceived := func (msg amqp.Delivery){
		log.Printf("Message received on exchange: %s\nWith routing key: %s\nWith body: %s", msg.Exchange, msg.RoutingKey, msg.Body)
	}

	estimateReceived := func(msg amqp.Delivery){
		genericMessageReceived(msg)
		data := msg.Body
		var estimation models.Estimation
		json.Unmarshal(data, &estimation)
		estimate, _ := services.CalculateCost(estimation)
		utils.PublishMessage("notification.exchange.notification", "notification.trip.estimatecalculated", estimate)
		log.Printf("Message received and unmarshaled into an estimation object: %s", estimation)
	}

	log.Printf("After defining generic MessageReceived and serve")

	go func() {
		utils.InitializeConsumer("eventBusTrip", "fanout", "eventQueueTripCalculation", "", genericMessageReceived)
	}()

	log.Printf("After initialize eventBusTrip Consumer")
	go func() {
		utils.InitializeConsumer("trip.exchange.tripcalculation", "topic", "trip.queue.calculationservice.calculatecost", "trip.estimation.estimatecalculated", estimateReceived)
	}()
	log.Printf("After initialize trip.exchange.tripcalculation Consumer")

	<- forever
}

func checkEurekaService(eurekaUp bool) bool {
	duration := time.Duration(15) * time.Second
	time.Sleep(duration)
	url := "http://discoveryservice:8761/eureka/"
	log.Println("Sending request to Eureka, waiting for response...")
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)
	if responseErr != nil {
		log.Printf("No response from Eureka, retrying...")
		return false
	}
	if response.Status != "204 No Content" {
		log.Printf("Success, Eureka was found!")
		return true
	}
	return false
}
