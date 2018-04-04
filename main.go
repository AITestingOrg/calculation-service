package main
	
import (
	"github.com/AITestingOrg/calculation-service/controllers"
	"github.com/AITestingOrg/calculation-service/eureka"
	"log"
	"time"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/api/v1/cost", controllers.GetCost).Methods("POST")
	log.Println("Calculation service is running...")
	var eurekaUp bool = false
	for eurekaUp != true {
		eurekaUp = checkEurekaService(eurekaUp)
	}
	eureka.PostToEureka()
	go eureka.StartHeartbeat()
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func checkEurekaService(eurekaUp bool) bool {
  	duration := time.Duration(15)*time.Second
	  time.Sleep(duration)
	url := "http://discovery-service:8761/eureka/"

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)
	if responseErr != nil {
		log.Printf("Response error")
		panic(responseErr)
	}
	if response.Status != "204 No Content" {
		log.Printf("Success, Eureka was found")
		return true
	}
	return false
}