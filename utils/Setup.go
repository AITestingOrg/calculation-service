package utils

import (
	"github.com/AITestingOrg/calculation-service/interfaces"
	"github.com/AITestingOrg/calculation-service/eureka"
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

func ProgramSetup(amqpPublisher interfaces.PublisherInterface, apiHandlers []interfaces.ApiHandlerInterface, amqpConsumers []interfaces.ConsumerInterface) {
	forever := make(chan bool)

	eureka.InitializeEurekaConnection()

	go amqpPublisher.InitializePublisher()

	go func() {
		r := mux.NewRouter()
		for _, handler := range apiHandlers {
			handler.AddHandlerToRouter(r)
		}
		http.Handle("/", r)
		log.Fatal(http.ListenAndServe(":8080", nil))
	} ()

	//initialize all of the rabbitmq consumers

	for _, consumer := range amqpConsumers {
		go consumer.InitializeConsumer()
	}

	<-forever
}
