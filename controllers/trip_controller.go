package controllers

import (
	"github.com/AITestingOrg/calculation-service/interfaces"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type CalculationServiceController struct {
	Handlers []interfaces.ApiHandlerInterface
}

func (controller CalculationServiceController) InitializeEndpoint() {
	r := mux.NewRouter()
	for _, handler := range controller.Handlers {
		r.HandleFunc(handler.GetPath(), handler.Handle).Methods(handler.GetRequestType())
	}
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
