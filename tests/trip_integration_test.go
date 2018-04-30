// +build integration

package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"calculation-service/controllers"
	"calculation-service/eureka"
	"calculation-service/models"

	"github.com/stretchr/testify/assert"
)

//Wait for Eureka to stand up before requesting any information during an integration test
func TestResponseStatusAndCost(t *testing.T) {
	log.Printf("Waiting for Eureka")
	var eurekaUp bool = false
	eurekaUp = eureka.CheckEurekaService(eurekaUp)
	log.Printf("Eureka found, creating request")
	trip := &models.Trip{
		Origin:      "Miami",
		Destination: "Weston",
	}
	jsonTrip, _ := json.Marshal(trip)
	request, _ := http.NewRequest("POST", "/api/v1/cost", bytes.NewBuffer(jsonTrip))
	request.Header.Set("Content-Type", "application/json")
	log.Printf("Finished creating request")

	if request.Method != "POST" {
		t.Errorf("Expected 'POST; request, got '%s'", request.Method)
	}
	if request.URL.EscapedPath() != "/api/v1/cost" {
		t.Errorf("Expected request to ‘/api/v1/cost’, got ‘%s’", request.URL.EscapedPath())
	}

	log.Print("Handling Request...")
	rr := httptest.NewRecorder()
	log.Printf("Recording")

	handler := http.HandlerFunc(controllers.GetCost)
	handler.ServeHTTP(rr, request)

	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	cost := m["cost"]

	//Asserts
	assert.Equal(t, 200, rr.Code, "OK response is expected")
	assert.NotNil(t, cost, "Cost should not be nil")
}
