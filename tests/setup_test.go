package tests

import (
	"github.com/AITestingOrg/calculation-service/interfaces"
	"github.com/AITestingOrg/calculation-service/tests/mocks"
	"github.com/AITestingOrg/calculation-service/utils"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestProgramSetup(t *testing.T) {
	mockPublisher := new(mocks.PublisherInterface)
	mockHandler1 := new(mocks.ApiHandlerInterface)
	mockConsumer1 := new(mocks.ConsumerInterface)
	mockHandler2 := new(mocks.ApiHandlerInterface)
	mockConsumer2 := new(mocks.ConsumerInterface)
	called := make(chan bool)

	apiHandlers := []interfaces.ApiHandlerInterface{mockHandler1, mockHandler2}
	amqpConsumers := []interfaces.ConsumerInterface{mockConsumer1, mockConsumer2}

	mockPublisher.On("InitializePublisher").Run(func(args mock.Arguments) { called <- true }).Return(nil)
	mockConsumer1.On("InitializeConsumer").Run(func(args mock.Arguments) { called <- true }).Return(nil)
	mockHandler1.On("AddHandlerToRouter", mock.Anything).Run(func(args mock.Arguments) { called <- true }).Return(nil)
	mockConsumer2.On("InitializeConsumer").Run(func(args mock.Arguments) { called <- true }).Return(nil)
	mockHandler2.On("AddHandlerToRouter", mock.Anything).Run(func(args mock.Arguments) { called <- true }).Return(nil)

	go utils.ProgramSetup(mockPublisher, apiHandlers, amqpConsumers)

	<-called
	<-called
	<-called
	<-called
	<-called

	mockPublisher.AssertExpectations(t)
	mockConsumer1.AssertExpectations(t)
	mockHandler1.AssertExpectations(t)
	mockConsumer2.AssertExpectations(t)
	mockHandler2.AssertExpectations(t)

	close(called)
}
