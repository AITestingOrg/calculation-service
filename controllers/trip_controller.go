package controllers

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/AITestingOrg/calculation-service/interfaces"
)

type CalculationServiceController struct {
	Handlers []interfaces.ApiHandlerInterface
}

func (controller CalculationServiceController) InitializeEndpoint(){
	r := mux.NewRouter()
	for _, handler := range controller.Handlers {
		r.HandleFunc(handler.GetPath(),handler.Handle).Methods(handler.GetRequestType())
	}
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
