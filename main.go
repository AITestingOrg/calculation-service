package main
	
import (
	"github.com/AITestingOrg/calculation-service/controllers"
	"github.com/AITestingOrg/calculation-service/eureka"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/api/v1/cost", controllers.GetCost).Methods("POST")
	log.Println("Calculation service is running...")
	var eurekaUp bool = false
	for eurekaUp != true {
		eurekaUp = eureka.CheckEurekaService(eurekaUp)
	}
	eureka.PostToEureka()
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}