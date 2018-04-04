package tests

import (
	"net/http"
    "net/http/httptest"
	"testing"
	"github.com/AITestingOrg/calculation-service/controllers"
	"github.com/AITestingOrg/calculation-service/models"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"time"
)

func TestResponseStatusAndCost(t *testing.T) {
	trip := &models.Trip {
		Origin:"Miami",
		Destination:"Weston",
	}
	jsonTrip, _ := json.Marshal(trip)
	request, _ := http.NewRequest("POST", "/api/v1/cost", bytes.NewBuffer(jsonTrip))
	request.Header.Set("Content-Type", "application/json")
	log.Printf("Finished creating request")

	if request.Method != "POST" {
		t.Errorf("Expected 'POST; request, got '%s'", request.Method)
	}
	if request.URL.EscapedPath() != "/api/v1/cost" {
		t.Errorf("Expected request to ‘/cost’, got ‘%s’", request.URL.EscapedPath())
	}
	
	log.Print("Handling Request...")
	rr := httptest.NewRecorder()
	log.Printf("Recording")

	duration := time.Duration(100)*time.Second
	time.Sleep(duration)
	handler := http.HandlerFunc(controllers.GetCost)
	handler.ServeHTTP(rr, request)

	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	cost := m["cost"]

	//Asserts
	assert.Equal(t, 200, rr.Code, "OK response is expected")
	assert.NotNil(t, cost, "Cost should not be nil")
}