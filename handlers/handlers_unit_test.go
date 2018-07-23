package handlers

import (
	"encoding/json"
	"github.com/AITestingOrg/calculation-service/models"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCalculateCost_HappyCase(t *testing.T) {
	currentDate := time.Now().Format("2006-01-02 03:04:05")
	gmapsEstimation := models.Estimation{
		Distance: 46180,
		Duration: 2552,
		Origin: "2250 N Commerce Pkwy, Weston, FL 33326",
		Destination: "3400 NE 163rd St, North Miami Beach, FL 33160",
		LastUpdated: currentDate,
	}

	cost := calculateCost(gmapsEstimation)
	assert.Equal(t,"2250 N Commerce Pkwy, Weston, FL 33326", cost.Origin)
	assert.Equal(t,"3400 NE 163rd St, North Miami Beach, FL 33160", cost.Destination)
	assert.Equal(t, 28.7, cost.Distance)
	assert.Equal(t, int64(42), cost.Duration)
	assert.Equal(t, 32.12, cost.Cost)
	assert.Equal(t, currentDate, cost.LastUpdated)
}

func TestEstimateReceived_HappyCase(t *testing.T) {
	currentDate := time.Now().Format("2006-01-02 03:04:05")
	gmapsEstimation := models.Estimation{
		Distance: 46180,
		Duration: 2552,
		Origin: "2250 N Commerce Pkwy, Weston, FL 33326",
		Destination: "3400 NE 163rd St, North Miami Beach, FL 33160",
		LastUpdated: currentDate,
	}
	estimationBody, _ := json.Marshal(gmapsEstimation)
	msg := amqp.Delivery{Exchange: "HappyExchange", RoutingKey: "HappyKey", Body: estimationBody}

	EstimateReceived(msg)
}
