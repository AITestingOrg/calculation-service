// +build !contract

package tests

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"testing"

// 	"github.com/pact-foundation/pact-go/dsl"
// 	"github.com/pact-foundation/pact-go/types"
// )

// func startServer() {
// 	mux := http.NewServeMux()
// 	var estimationReturned string

// 	mux.HandleFunc("/api/v1/cost", func(w http.ResponseWriter, req *http.Request) {
// 		w.Header().Add("Content-Type", "application/json")
// 		fmt.Fprintf(w, ` { "originAddress": "",
// 			"destinationAddress": "",
// 			"distance": 0.0,
// 			"duration": 0,
// 			"cost": 0.0,
// 			"lastUpdated": "" } `)
// 	})

// 	mux.HandleFunc("/setup", func(w http.ResponseWriter, req *http.Request) {
// 		var s *types.ProviderState
// 		decoder := json.NewDecoder(req.Body)
// 		decoder.Decode(&s)
// 		if s.State == "STATE TO BE GIVEN BY TRIP-CMD" {
// 			estimationReturned =
// 				` { "originAddress": ` + dsl.Like("9700 Collins Ave, Bal Harbour, FL 33154") + `,
// 				"destinationAddress": ` + dsl.Like("2250 N Commerce Pkwy, Weston, FL 33326") + `,
// 				"distance": ` + dsl.Like("2000.0") + `,
// 				"duration": ` + dsl.Like("2000") + `,
// 				"cost": ` + dsl.Like("2000.0") + `,
// 				"distance": ` + dsl.Like("Apr 09 2018") + `, } `
// 		} else {
// 			panic("Provider state is invalid")
// 		}
// 		w.Header().Add("Content-Type", "application/json")
// 	})

// 	go http.ListenAndServe(":8080", mux)
// }

// //Calculation provides for Trip Management Service
// func TestProvider(t *testing.T) {
// 	pact := &dsl.Pact{
// 		Port:     6666, // Ensure this port matches the daemon port!
// 		Consumer: "Trip Management Service",
// 		Provider: "Calculation Service",
// 	}

// 	// Start provider API in the background
// 	go startServer()
// 	var dir, _ = os.Getwd()
// 	var pactDir = fmt.Sprintf("%s/pacts", dir)
// 	// Verify the Provider with local Pact Files
// 	pact.VerifyProvider(t, types.VerifyRequest{
// 		ProviderBaseURL:        "http://localhost:8080",
// 		PactURLs:               []string{filepath.ToSlash(fmt.Sprintf("%s/PATH TO PACT JSON", pactDir))},
// 		ProviderStatesSetupURL: "http://localhost:8080/setup",
// 	})
// }
