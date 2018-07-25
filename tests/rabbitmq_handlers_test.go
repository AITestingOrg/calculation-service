package tests

import (
	"testing"
	"github.com/AITestingOrg/calculation-service/handlers"
	"github.com/AITestingOrg/calculation-service/models"
	"github.com/streadway/amqp"
	"encoding/json"
	"github.com/AITestingOrg/calculation-service/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
	"errors"
)

func TestEstimateHandler_HappyCase(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.EstimateHandler{Publisher: mockPublisher}

	lastUpdated := time.Now().Format("2006-01-02 03:04:05")
	happyEstimate := models.Estimation{Origin: "Weston, Fl", Destination: "Miami, Fl", UserId: "6e012898-8a5b-4959-92d6-8b7d669384b4", Duration: 3000, Distance: 3000}
	var realDuration = float64(happyEstimate.Duration/ 60)
	var realDistance = float64(int(happyEstimate.Distance/1609*100)) / 100
	var cost = float64(int(((realDuration * 0.15)+(realDistance * 0.9))*100)) / 100
	happyEstimateByteArray, _ := json.Marshal(happyEstimate)

	estimationMatcher := mock.MatchedBy(func(estimation models.Estimation) bool {
		res := estimation.Origin == happyEstimate.Origin &&
			estimation.Destination == happyEstimate.Destination &&
			string(estimation.LastUpdated)[0:15] == string(lastUpdated)[0:15] &&
			estimation.UserId == happyEstimate.UserId &&
			estimation.Distance == realDistance &&
			estimation.Duration == int64(realDuration) &&
			estimation.Cost == cost
		return res
		})

	mockPublisher.Mock.On("PublishMessage", "notification.exchange.notification", "notification.trip.estimatecalculated", estimationMatcher).Return(nil)

	//Act
	err := handler.Handle(amqp.Delivery{Body: happyEstimateByteArray})

	//Assert
	assert.Equal(t, err, nil)
}

func TestEstimateHandler_ZeroDistance(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.EstimateHandler{Publisher: mockPublisher}

	lastUpdated := time.Now().Format("2006-01-02 03:04:05")
	happyEstimate := models.Estimation{Origin: "Weston, Fl", Destination: "Miami, Fl", UserId: "6e012898-8a5b-4959-92d6-8b7d669384b4", Duration: 3000, Distance: 0}
	var realDuration = float64(happyEstimate.Duration/ 60)
	var realDistance = float64(int(happyEstimate.Distance/1609*100)) / 100
	var cost = float64(int(((realDuration * 0.15)+(realDistance * 0.9))*100)) / 100
	happyEstimateByteArray, _ := json.Marshal(happyEstimate)

	estimationMatcher := mock.MatchedBy(func(estimation models.Estimation) bool {
		res := estimation.Origin == happyEstimate.Origin &&
			estimation.Destination == happyEstimate.Destination &&
			string(estimation.LastUpdated)[0:15] == string(lastUpdated)[0:15] &&
			estimation.UserId == happyEstimate.UserId &&
			estimation.Distance == realDistance &&
			estimation.Duration == int64(realDuration) &&
			estimation.Cost == cost
		return res
	})

	mockPublisher.Mock.On("PublishMessage", "notification.exchange.notification", "notification.trip.estimatecalculated", estimationMatcher).Return(nil)

	//Act
	err := handler.Handle(amqp.Delivery{Body: happyEstimateByteArray})

	//Assert
	assert.Equal(t, err, nil)
}

func TestEstimateHandler_ZeroDuration(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.EstimateHandler{Publisher: mockPublisher}

	lastUpdated := time.Now().Format("2006-01-02 03:04:05")
	happyEstimate := models.Estimation{Origin: "Weston, Fl", Destination: "Miami, Fl", UserId: "6e012898-8a5b-4959-92d6-8b7d669384b4", Duration: 0, Distance: 3000}
	var realDuration = float64(happyEstimate.Duration/ 60)
	var realDistance = float64(int(happyEstimate.Distance/1609*100)) / 100
	var cost = float64(int(((realDuration * 0.15)+(realDistance * 0.9))*100)) / 100
	happyEstimateByteArray, _ := json.Marshal(happyEstimate)

	estimationMatcher := mock.MatchedBy(func(estimation models.Estimation) bool {
		res := estimation.Origin == happyEstimate.Origin &&
			estimation.Destination == happyEstimate.Destination &&
			string(estimation.LastUpdated)[0:15] == string(lastUpdated)[0:15] &&
			estimation.UserId == happyEstimate.UserId &&
			estimation.Distance == realDistance &&
			estimation.Duration == int64(realDuration) &&
			estimation.Cost == cost
		return res
	})

	mockPublisher.Mock.On("PublishMessage", "notification.exchange.notification", "notification.trip.estimatecalculated", estimationMatcher).Return(nil)

	//Act
	err := handler.Handle(amqp.Delivery{Body: happyEstimateByteArray})

	//Assert
	assert.Equal(t, err, nil)
}

func TestEstimateHandler_PublisherFailed(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.EstimateHandler{Publisher: mockPublisher}

	lastUpdated := time.Now().Format("2006-01-02 03:04:05")
	happyEstimate := models.Estimation{Origin: "Weston, Fl", Destination: "Miami, Fl", UserId: "6e012898-8a5b-4959-92d6-8b7d669384b4", Duration: 3000, Distance: 3000}
	var realDuration = float64(happyEstimate.Duration/ 60)
	var realDistance = float64(int(happyEstimate.Distance/1609*100)) / 100
	var cost = float64(int(((realDuration * 0.15)+(realDistance * 0.9))*100)) / 100
	happyEstimateByteArray, _ := json.Marshal(happyEstimate)

	estimationMatcher := mock.MatchedBy(func(estimation models.Estimation) bool {
		res := estimation.Origin == happyEstimate.Origin &&
			estimation.Destination == happyEstimate.Destination &&
			string(estimation.LastUpdated)[0:15] == string(lastUpdated)[0:15] &&
			estimation.UserId == happyEstimate.UserId &&
			estimation.Distance == realDistance &&
			estimation.Duration == int64(realDuration) &&
			estimation.Cost == cost
		return res
	})

	mockErr := errors.New("error Message")
	mockPublisher.Mock.On("PublishMessage", "notification.exchange.notification", "notification.trip.estimatecalculated", estimationMatcher).Return(mockErr)
	mockErr = errors.New("error publishing message: " + mockErr.Error())

	//Act
	err := handler.Handle(amqp.Delivery{Body: happyEstimateByteArray})

	//Assert
	assert.Equal(t, mockErr, err)
}

func TestEstimateHandler_UnmarshalError(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.EstimateHandler{Publisher: mockPublisher}

	invalidEstimate := []byte("byte")
	var estimate models.Estimation
	expectedErr := json.Unmarshal(invalidEstimate, &estimate)

	//Act
	err := handler.Handle(amqp.Delivery{Body: invalidEstimate})

	//Assert
	assert.Equal(t, errors.New("error unmarshalling data into an estimation object: " + expectedErr.Error()), err)
}

func TestEstimateHandler_InvalidOrigin(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.EstimateHandler{Publisher: mockPublisher}

	estimate := models.Estimation{Origin: "", Destination: "Miami, Fl", UserId: "6e012898-8a5b-4959-92d6-8b7d669384b4", Duration: 3000, Distance: 3000}
	estimateByteArray, _ := json.Marshal(estimate)

	expectedErr := estimate.ValidateFields("originAddress", "destinationAddress", "distance", "duration", "userId")

	//Act
	err := handler.Handle(amqp.Delivery{Body: estimateByteArray})

	//Assert
	assert.Equal(t, errors.New("error with the parsed estimation object: \n" + expectedErr.Error()), err)
}

func TestEstimateHandler_InvalidDestination(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.EstimateHandler{Publisher: mockPublisher}

	estimate := models.Estimation{Origin: "Weston, Fl", Destination: "", UserId: "6e012898-8a5b-4959-92d6-8b7d669384b4", Duration: 3000, Distance: 3000}
	estimateByteArray, _ := json.Marshal(estimate)

	expectedErr := estimate.ValidateFields("originAddress", "destinationAddress", "distance", "duration", "userId")

	//Act
	err := handler.Handle(amqp.Delivery{Body: estimateByteArray})

	//Assert
	assert.Equal(t, errors.New("error with the parsed estimation object: \n" + expectedErr.Error()), err)
}

func TestEstimateHandler_InvalidDistance(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.EstimateHandler{Publisher: mockPublisher}

	estimate := models.Estimation{Origin: "Weston, Fl", Destination: "Miami, Fl", UserId: "6e012898-8a5b-4959-92d6-8b7d669384b4", Duration: 3000, Distance: -1}
	estimateByteArray, _ := json.Marshal(estimate)

	expectedErr := estimate.ValidateFields("originAddress", "destinationAddress", "distance", "duration", "userId")

	//Act
	err := handler.Handle(amqp.Delivery{Body: estimateByteArray})

	//Assert
	assert.Equal(t, errors.New("error with the parsed estimation object: \n" + expectedErr.Error()), err)
}

func TestEstimateHandler_InvalidDuration(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.EstimateHandler{Publisher: mockPublisher}

	estimate := models.Estimation{Origin: "Weston, Fl", Destination: "Miami, Fl", UserId: "6e012898-8a5b-4959-92d6-8b7d669384b4", Duration: -1, Distance: 3000}
	estimateByteArray, _ := json.Marshal(estimate)

	expectedErr := estimate.ValidateFields("originAddress", "destinationAddress", "distance", "duration", "userId")

	//Act
	err := handler.Handle(amqp.Delivery{Body: estimateByteArray})

	//Assert
	assert.Equal(t, errors.New("error with the parsed estimation object: \n" + expectedErr.Error()), err)
}

func TestEstimateHandler_InvalidUUID(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.EstimateHandler{Publisher: mockPublisher}

	estimate := models.Estimation{Origin: "Weston, Fl", Destination: "Miami, Fl", UserId: "0", Duration: 3000, Distance: 3000}
	estimateByteArray, _ := json.Marshal(estimate)

	expectedErr := estimate.ValidateFields("originAddress", "destinationAddress", "distance", "duration", "userId")

	//Act
	err := handler.Handle(amqp.Delivery{Body: estimateByteArray})

	//Assert
	assert.Equal(t, errors.New("error with the parsed estimation object: \n" + expectedErr.Error()), err)
}
