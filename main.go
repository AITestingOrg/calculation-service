package main
	
import (
	"calculation-service/controllers"
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"calculation-service/eureka"
	"time"
)

func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/cost", controllers.GetCost).Methods("POST")
	log.Println("Calculation service is running...")

	duration := time.Duration(20)*time.Second
	time.Sleep(duration)
	// var eurekaUp bool = false
	// for eurekaUp != true {
	// 	eurekaUp = checkEurekaService(eurekaUp)
	// }
	eureka.PostToEureka()
	log.Fatal(http.ListenAndServe(":8000", r))
}

func checkEurekaService(eurekaUp bool) bool {
  	duration := time.Duration(20)*time.Second
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