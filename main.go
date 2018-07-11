package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AITestingOrg/calculation-service/controllers"
	"github.com/AITestingOrg/calculation-service/eureka"

	"github.com/gorilla/mux"

	"fmt"
	"github.com/streadway/amqp"
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

	//http.Handle("/", r)
	//log.Fatal(http.ListenAndServe(":8080", nil))


	log.Printf("before dialing rabbitmq")

	conn, err := amqp.Dial("amqp://guest:guest@" + os.Getenv("RABBIT_HOST") + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	log.Printf("after dialing rabbitmq")

	err = ch.ExchangeDeclare(
		"trip-calculation-exchange",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare an exchange")

	log.Printf("after declaring exchange")

	q, err := ch.QueueDeclare(
		"trip-calculation-queue",
		false,
		false,
		false,
		false,
		nil,
	)



	log.Printf("Declared queue: " + q.Name)

	failOnError(err, "Failed to declare the queue, trip-calculation-queue")

	body := "hello from calculation service"
	err = ch.Publish(
		"trip-calculation-exchange",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: 	"text/plain",
			Body:			[]byte(body),
		})
	failOnError(err, "Failed to publish the hello message")

	log.Printf("Declared queue: " + q.Name)

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages")
	<-forever
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
