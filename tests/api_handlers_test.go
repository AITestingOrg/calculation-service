package tests

import (
	"bytes"
	"encoding/json"
	"github.com/AITestingOrg/calculation-service/handlers"
	"github.com/AITestingOrg/calculation-service/models"
	"github.com/AITestingOrg/calculation-service/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestCostEstimateHandler_HappyCase(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.CostEstimateHandler{Publisher: mockPublisher}
	trip := models.Trip{Origin: "Miami, Fl", Destination: "Weston, Fl", UserId: "6dc35c49-0e20-4394-8fa1-2532a830067d"}
	tripBytes, _ := json.Marshal(trip)
	body := bytes.NewBuffer(tripBytes)
	mockRequest := httptest.NewRequest("POST", "/api/v1/cost", body)
	w := httptest.NewRecorder()

	tripMatcher := mock.MatchedBy(func(tripCheck models.Trip) bool {
		return tripCheck.Origin == trip.Origin &&
			tripCheck.Destination == trip.Destination &&
			tripCheck.UserId == trip.UserId
	})

	mockPublisher.Mock.On("PublishMessage", "trip.exchange.tripcalculation", "trip.estimation.estimaterequested", tripMatcher).Return(nil)

	//Act
	handler.Handle(w, mockRequest)

	resp := w.Result()
	respBody, _ := ioutil.ReadAll(resp.Body)

	//Assert
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, resp.Header.Get("Content-Type"), "text/plain")
	assert.Equal(t, string(respBody), "Trip Estimate Request Retrieved. Forwarding request to Gmaps Adapter")
	mockPublisher.Mock.AssertExpectations(t)
}

func TestCostEstimateHandler_EmptyOrigin(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.CostEstimateHandler{Publisher: mockPublisher}
	trip := models.Trip{Origin: "", Destination: "Weston, Fl", UserId: "6dc35c49-0e20-4394-8fa1-2532a830067d"}
	tripBytes, _ := json.Marshal(trip)
	body := bytes.NewBuffer(tripBytes)
	mockRequest := httptest.NewRequest("POST", "/api/v1/cost", body)
	w := httptest.NewRecorder()

	//Act
	handler.Handle(w, mockRequest)

	resp := w.Result()
	respBody, _ := ioutil.ReadAll(resp.Body)

	//Assert
	assert.Equal(t, 400, resp.StatusCode)
	expectedErr := trip.ValidateFields("originAddress", "destinationAddress", "userId")
	assert.Equal(t, "ERROR: Invalid trip arguments:\n"+expectedErr.Error()+"\n", string(respBody))
	mockPublisher.Mock.AssertNotCalled(t, "PublishMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestCostEstimateHandler_EmptyDestination(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.CostEstimateHandler{Publisher: mockPublisher}
	trip := models.Trip{Origin: "Miami, Fl", Destination: "", UserId: "6dc35c49-0e20-4394-8fa1-2532a830067d"}
	tripBytes, _ := json.Marshal(trip)
	body := bytes.NewBuffer(tripBytes)
	mockRequest := httptest.NewRequest("POST", "/api/v1/cost", body)
	w := httptest.NewRecorder()

	//Act
	handler.Handle(w, mockRequest)

	resp := w.Result()
	respBody, _ := ioutil.ReadAll(resp.Body)

	//Assert
	assert.Equal(t, 400, resp.StatusCode)
	expectedErr := trip.ValidateFields("originAddress", "destinationAddress", "userId")
	assert.Equal(t, "ERROR: Invalid trip arguments:\n"+expectedErr.Error()+"\n", string(respBody))
	mockPublisher.Mock.AssertNotCalled(t, "PublishMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestCostEstimateHandler_InvalidUserId(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.CostEstimateHandler{Publisher: mockPublisher}
	trip := models.Trip{Origin: "Miami, Fl", Destination: "Weston, Fl", UserId: ""}
	tripBytes, _ := json.Marshal(trip)
	body := bytes.NewBuffer(tripBytes)
	mockRequest := httptest.NewRequest("POST", "/api/v1/cost", body)
	w := httptest.NewRecorder()

	//Act
	handler.Handle(w, mockRequest)

	resp := w.Result()
	respBody, _ := ioutil.ReadAll(resp.Body)

	//Assert
	assert.Equal(t, 400, resp.StatusCode)
	expectedErr := trip.ValidateFields("originAddress", "destinationAddress", "userId")
	assert.Equal(t, "ERROR: Invalid trip arguments:\n"+expectedErr.Error()+"\n", string(respBody))
	mockPublisher.Mock.AssertNotCalled(t, "PublishMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestCostEstimateHandler_NonExistentOrigin(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.CostEstimateHandler{Publisher: mockPublisher}
	tripBytes := bytes.NewBuffer([]byte("{\"destination\":\"Miami, Fl\",\"userId\":\"6dc35c49-0e20-4394-8fa1-2532a830067d\"}"))
	mockRequest := httptest.NewRequest("POST", "/api/v1/cost", tripBytes)
	w := httptest.NewRecorder()

	//Act
	handler.Handle(w, mockRequest)

	resp := w.Result()
	respBody, _ := ioutil.ReadAll(resp.Body)

	//Assert
	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "ERROR: Invalid trip arguments:\nInvalid originAddress.\n\tGiven: \n\tExpected: Non Empty String\n\n", string(respBody))
	mockPublisher.Mock.AssertNotCalled(t, "PublishMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestCostEstimateHandler_NonExistentDestination(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.CostEstimateHandler{Publisher: mockPublisher}
	tripBytes := bytes.NewBuffer([]byte("{\"origin\":\"Weston, Fl\",\"userId\":\"6dc35c49-0e20-4394-8fa1-2532a830067d\"}"))
	mockRequest := httptest.NewRequest("POST", "/api/v1/cost", tripBytes)
	w := httptest.NewRecorder()

	//Act
	handler.Handle(w, mockRequest)

	resp := w.Result()
	respBody, _ := ioutil.ReadAll(resp.Body)

	//Assert
	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "ERROR: Invalid trip arguments:\nInvalid destinationAddress.\n\tGiven: \n\tExpected: Non Empty String\n\n", string(respBody))
	mockPublisher.Mock.AssertNotCalled(t, "PublishMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestCostEstimateHandler_NonExistentUserId(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.CostEstimateHandler{Publisher: mockPublisher}
	tripBytes := bytes.NewBuffer([]byte("{\"origin\":\"Weston, Fl\",\"destination\":\"Miami, Fl\"}"))
	mockRequest := httptest.NewRequest("POST", "/api/v1/cost", tripBytes)
	w := httptest.NewRecorder()

	//Act
	handler.Handle(w, mockRequest)

	resp := w.Result()
	respBody, _ := ioutil.ReadAll(resp.Body)

	//Assert
	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "ERROR: Invalid trip arguments:\nInvalid userId.\n\tGiven: \n\tExpected: Non-empty UUID\n\n", string(respBody))
	mockPublisher.Mock.AssertNotCalled(t, "PublishMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestCostEstimateHandler_ExtraArgumentsProvided(t *testing.T) {
	//Arrange
	mockPublisher := new(mocks.PublisherInterface)

	handler := handlers.CostEstimateHandler{Publisher: mockPublisher}
	tripBytes := bytes.NewBuffer([]byte("{\"origin\":\"Weston, Fl\",\"destination\":\"Miami, Fl\",\"userId\":\"6dc35c49-0e20-4394-8fa1-2532a830067d\",\"departureTime\":-1}"))
	mockRequest := httptest.NewRequest("POST", "/api/v1/cost", tripBytes)
	w := httptest.NewRecorder()

	tripMatcher := mock.MatchedBy(func(tripCheck models.Trip) bool {
		return tripCheck.Origin == "Weston, Fl" &&
			tripCheck.Destination == "Miami, Fl" &&
			tripCheck.UserId == "6dc35c49-0e20-4394-8fa1-2532a830067d" &&
			tripCheck.DepartureTime != -1
	})

	mockPublisher.Mock.On("PublishMessage", "trip.exchange.tripcalculation", "trip.estimation.estimaterequested", tripMatcher).Return(nil)

	//Act
	handler.Handle(w, mockRequest)

	resp := w.Result()
	respBody, _ := ioutil.ReadAll(resp.Body)

	//Assert
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, resp.Header.Get("Content-Type"), "text/plain")
	assert.Equal(t, string(respBody), "Trip Estimate Request Retrieved. Forwarding request to Gmaps Adapter")
	mockPublisher.Mock.AssertExpectations(t)
}
