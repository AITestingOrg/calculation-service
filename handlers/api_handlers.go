package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/AITestingOrg/calculation-service/db"
	"github.com/AITestingOrg/calculation-service/interfaces"
	"github.com/AITestingOrg/calculation-service/models"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type CostEstimateHandler struct {
	Publisher interfaces.PublisherInterface
}

type CostStorage struct {
	Handler interfaces.ControllerInterface
}

func (handler CostEstimateHandler) AddHandlerToRouter(r *mux.Router) {
	r.HandleFunc("/api/v1/cost", handler.Handle).Methods("POST")
}

func (handler CostStorage) AddHandlerToRouter(r *mux.Router) {
	r.HandleFunc("/api/v1/cost", handler.Handle).Methods("GET")
}

func (handler CostStorage) Handle(w http.ResponseWriter, r *http.Request) {
	session, err := db.Connect()
	defer session.Close()
	if session == nil {
		log.Fatalf("No Session")
		log.Println("session is Nil", err)
	}
	log.Println("Connecting to TRIPCOST database...")
	c := session.DB("TRIPCOST").C("costs")
	cost := models.Cost{}
	errJson := json.NewDecoder(r.Body).Decode(&cost)
	if errJson != nil {
		panic(errJson)
	}

	log.Println("GET request UUID retreived: ", cost.UserId)
	err = c.Find(bson.M{"userId": cost.UserId}).Limit(1).One(&cost)
	if err != nil {
		err = errors.New("Failed to find cost: " + err.Error())
		log.Print(err)
		http.Error(w, err.Error(), 404)
		return
	} else {
		log.Printf("Found cost")
	}
	return
}

func (handler CostEstimateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	//Converts objects into json
	var trip models.Trip
	json.Unmarshal(body, &trip)

	log.Printf("Validating trip and estimation body...")
	err := trip.ValidateFields("originAddress", "destinationAddress", "userId")
	if err != nil {
		err = errors.New("ERROR: Invalid trip arguments:\n" + err.Error())
		log.Print(err)
		http.Error(w, err.Error(), 400)
		return
	}

	//Sets departure time to current time
	log.Printf("Setting current time for TripEstimate...")
	currentTime := time.Now().Unix()
	trip.DepartureTime = currentTime

	//Emit to trip.exchange.tripcalculation, with key trip.estimation.estimaterequested
	err = handler.Publisher.PublishMessage("trip.exchange.tripcalculation", "trip.estimation.estimaterequested", trip)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(500)
		w.Write([]byte("Error publishing message to rabbitmq: " + err.Error()))
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Trip Estimate Request Retrieved. Forwarding request to Gmaps Adapter"))
}
