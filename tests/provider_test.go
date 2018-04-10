package tests

import (
	"os"
	"fmt"
	"log"
	"testing"
	"net/http"
	"path/filepath"
	"encoding/json"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/AITestingOrg/calculation-service/models"
)

func startServer() {
	mux := http.NewServeMux()
	var estimation models.Estimation

	mux.HandleFunc("/api/v1/cost", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, ` { "originAddress": "",
			"destinationAddress": "",
			"distance": 0.0,
			"duration": 0,
			"cost": 0.0,
			"lastUpdated": "" } `)
	})

	mux.HandleFunc("/setup", func(w http.ResponseWriter, req *http.Request) {
		var s *types.ProviderState
		decoder := json.NewDecoder(req.Body)
		decoder.Decode(&s)
		if s.State == "Origin, destination, and current time" {
			estimation = models.Estimation {
				Origin: "9700 Collins Ave, Bal Harbour, FL 33154", 
				Destination: "2250 N Commerce Pkwy, Weston, FL 33326", 
				Distance: 32.12,
				Duration: 40,
				Cost: 34.9,
				LastUpdated: "Apr 09 2018",
			}
		} else {
			panic("Provider state is invalid")
		}
		w.Header().Add("Content-Type", "application/json")
	})

	go http.ListenAndServe(":8000", mux)
}

//Calculation provides for Trip Management Service
func TestProvider(t *testing.T) {
	pact := &dsl.Pact{
		Port:     6666, // Ensure this port matches the daemon port!
		Consumer: "Trip Management Service",
		Provider: "Calculation Service",
	}
  
	// Start provider API in the background
	go startServer()
	log.Printf("Server started")
	var dir, _ = os.Getwd()
	var pactDir = fmt.Sprintf("%s/pacts", dir)
	// Verify the Provider with local Pact Files
	pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL:        "http://localhost:8000",
		PactURLs:               []string{filepath.ToSlash(fmt.Sprintf("%s/trip_management_service-calculation_service.json", pactDir))},
		ProviderStatesSetupURL: "http://localhost:8000/setup",
	})
}