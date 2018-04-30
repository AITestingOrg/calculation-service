// +build !unit !integration

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"calculation-service/models"

	"github.com/pact-foundation/pact-go/dsl"
)

//Calculation consumes from the Gmaps adapter
func TestConsumer(t *testing.T) {
	// Create Pact connecting to local Daemon
	pact := &dsl.Pact{
		Port:     6666, // Ensure this port matches the daemon port!
		Consumer: "Calculation Service",
		Provider: "Gmaps Adapter",
		Host:     "localhost",
	}

	defer pact.Teardown()

	//Trip Model
	trip := models.Trip{
		Origin:        "9700 Collins Ave, Bal Harbour, FL 33154",
		Destination:   "2250 N Commerce Pkwy, Weston, FL 33326",
		DepartureTime: 1523340999999999,
	}

	// Pass in test case
	var test = func() error {
		url := fmt.Sprintf("http://localhost:%d/api/v1/directions", pact.Server.Port)
		encodedTrip, marshallErr := json.Marshal(trip)
		if marshallErr != nil {
			log.Printf("Error encoding trip")
			panic(marshallErr)
		}
		request, err := http.NewRequest("POST", url, bytes.NewBuffer(encodedTrip))
		if err != nil {
			log.Printf("Error creating new POST request")
			panic(err)
		}
		request.Header.Set("Content-Type", "application/json")

		if _, err = http.DefaultClient.Do(request); err != nil {
			log.Printf("Error executing POST request to gmaps-adapter")
			panic(err)
		}
		return err
	}

	// Set up our expected interactions.
	pact.
		AddInteraction().
		Given("Origin, destination, and current time").
		UponReceiving("a POST request to get the duration and distance").
		WithRequest(dsl.Request{
			Method:  "POST",
			Path:    "/api/v1/directions",
			Headers: map[string]string{"Content-Type": "application/json"},
			Body:    trip,
		}).
		WillRespondWith(dsl.Response{
			Status:  200,
			Headers: map[string]string{"Content-Type": "application/json"},
			Body: `{ "distance": ` + dsl.Like("2000") + `,
				"duration": ` + dsl.Like("2000") + `
				}`,
		})

	// Verify
	if err := pact.Verify(test); err != nil {
		log.Fatalf("Error on Verify: %v", err)
	}
}
