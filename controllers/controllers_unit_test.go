// +build unit

package controllers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/AITestingOrg/calculation-service/models"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	http.HandleFunc("/api/v1/cost", GetCost)
}

func TestGetCost(t *testing.T) {

	data :=
		`<instance>
		</instance>`

	estimation := &models.Estimation{
		Distance: 59485,
		Duration: 4736,
	}

	encodedEstimation, err := json.Marshal(estimation)
	if err != nil {
		panic(err)
	}

	encodedData, err := xml.Marshal(data)
	if err != nil {
		panic(err)
	}

	httpmock.Activate()
	httpmock.RegisterResponder("GET", "http://discovery-service:8761/eureka/apps/gmapsadapter", httpmock.NewBytesResponder(200, encodedData))
	httpmock.RegisterResponder("POST", "http://localhost:8080/api/v1/directions", httpmock.NewBytesResponder(200, encodedEstimation))
	defer httpmock.Deactivate()

	trip := &models.Trip{
		Origin:      "Miami",
		Destination: "Weston",
	}

	jsonTrip, _ := json.Marshal(trip)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/cost", bytes.NewBuffer(jsonTrip))
	response := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(response, request)

	json.NewDecoder(response.Body).Decode(&estimation)

	assert.Equal(t, 200, response.Code)
	assert.NotNil(t, response.Body)
	assert.Equal(t, "Miami", estimation.Origin)
	assert.Equal(t, "Weston", estimation.Destination)
	assert.Equal(t, 36.97, estimation.Distance)
	assert.Equal(t, int64(78), estimation.Duration)
	assert.Equal(t, 44.97, estimation.Cost)

}
