package handlers

import (
	"encoding/json"
	"errors"
	"github.com/AITestingOrg/calculation-service/interfaces"
	"github.com/AITestingOrg/calculation-service/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"github.com/AITestingOrg/calculation-service/db"
	"strings"
	"gopkg.in/mgo.v2/bson"
)

type CostEstimateHandler struct {
	Publisher interfaces.PublisherInterface
}
type CostStorage struct {
	Publisher interfaces.PublisherInterface
}

func (handler CostEstimateHandler) AddHandlerToRouter(r *mux.Router) {
	r.HandleFunc("/api/v1/cost", handler.Handle).Methods("POST")
}

func (handler CostStorage) AddHandlerToRouter(r *mux.Router) {
	r.HandleFunc("/api/v1/cost/{id}", handler.Handle).Methods("GET")
}

func (handler CostStorage) Handle(w http.ResponseWriter, r *http.Request) {
	session := db.MgoSession.Copy()
	defer session.Close()

	if session == nil {
		log.Println("session is Nil")
	}

	c := session.DB("TRIPCOST").C("costs")
	param := strings.TrimPrefix(r.URL.Path, "/api/v1/cost/")
	log.Println("GET parameter retreived: ", param)

	var cost models.Cost
	err := c.Find(bson.M{"userId": param}).Sort("-departureTime").Limit(1).One(&cost)

	if err != nil {
		log.Printf("Failed to find cost")
	}
	log.Printf("Found cost")
	log.Println(cost.Origin, cost.Destination, cost.Cost, cost.DepartureTime)
	return
}

func (handler CostEstimateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
//	//Converts objects into json
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
