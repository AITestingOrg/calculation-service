// +build unit

package services

import (
	"encoding/json"
	"encoding/xml"
	"github.com/AITestingOrg/calculation-service/models"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetGmapsEstimation(t *testing.T) {

	data :=
		`<instance>
			<ipAddr>localhost</ipAddr>
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

	result := GetGmapsEstimation(*trip)

	assert.NotNil(t, result)
	assert.Equal(t, float64(59485), result.Distance)
	assert.Equal(t, int64(4736), result.Duration)

}

func TestGetGmapsEstimationWhenConnectionRefused(t *testing.T) {

	trip := &models.Trip{
		Origin:      "Narnia",
		Destination: "Chicago",
	}

	assert.Panics(t, func() { GetGmapsEstimation(*trip) }, "The code did not panic")

}

func TestCalculateCost(t *testing.T) {

	data :=
		`<instance>
			<ipAddr>localhost</ipAddr>
		</instance>`

	estimation := &models.Estimation{
		Distance: 46180,
		Duration: 2552,
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
		Origin:      "2250 N Commerce Pkwy, Weston, FL 33326",
		Destination: "3400 NE 163rd St, North Miami Beach, FL 33160",
	}

	result := CalculateCost(*trip, *estimation)
	resultBody := models.Estimation{}
	json.Unmarshal(result, &resultBody)

	assert.NotNil(t, result)
	assert.Equal(t, "2250 N Commerce Pkwy, Weston, FL 33326", resultBody.Origin)
	assert.Equal(t, "3400 NE 163rd St, North Miami Beach, FL 33160", resultBody.Destination)
	assert.Equal(t, 28.7, resultBody.Distance)
	assert.Equal(t, int64(42), resultBody.Duration)
	assert.Equal(t, 32.12, resultBody.Cost)
}

func TestCalculateCostWhenConnectionRefused(t *testing.T) {

	trip := &models.Trip{
		Origin:      "Narnia",
		Destination: "Chicago",
	}

	estimation := &models.Estimation{
		Distance: 21344,
		Duration: 231,
	}

	assert.Panics(t, func() { CalculateCost(*trip, *estimation) }, "The code did not panic")

}