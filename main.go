package main
	
import (
	"calculation-service/controllers"
	"net/http"
	"log"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/cost", controllers.GetCost).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}