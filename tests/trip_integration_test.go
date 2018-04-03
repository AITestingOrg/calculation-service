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
)

func TestResponseStatusAndCost(t *testing.T) {
	trip := &models.Trip {
		Origin:"Miami",
		Destination:"Weston",
	}
	jsonTrip, _ := json.Marshal(trip)
	request, _ := http.NewRequest("POST", "/cost", bytes.NewBuffer(jsonTrip))

	if request.Method != "POST" {
		t.Errorf("Expected 'POST; request, got '%s'", request.Method)
	}
	if request.URL.EscapedPath() != "/cost" {
		t.Errorf("Expected request to ‘/cost’, got ‘%s’", request.URL.EscapedPath())
	}
	
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.GetCost)
	handler.ServeHTTP(rr, request)

	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	cost := m["cost"]

	//Asserts
	assert.Equal(t, 200, rr.Code, "OK response is expected")
	assert.NotNil(t, cost, "Cost should not be nil")
}