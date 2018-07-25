package tests

import (
	"github.com/AITestingOrg/calculation-service/handlers"
	"github.com/AITestingOrg/calculation-service/tests/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildPublisher(t *testing.T) {
	// Arrange
	mockMain := new(mocks.MainProgramInterface)
	mainHandler := handlers.MainProgram{mockMain}

	// Act
	producer := mainHandler.BuildPublisher()

	// Assert
	assert.NotNil(t, producer)
}

func TestBuildConsumer(t *testing.T) {
	// Arrange
	mockMain := new(mocks.MainProgramInterface)
	realMain := new(handlers.MainProgram)
	mainHandler := handlers.MainProgram{mockMain}
	mockPublisher := mainHandler.BuildPublisher()

	// Act
	consumer := realMain.BuildConsumer(mockPublisher)

	// Assert
	assert.NotNil(t, consumer)
}

func TestBuildController(t *testing.T) {
	// Arrange
	mockMain := new(mocks.MainProgramInterface)
	realMain := new(handlers.MainProgram)
	mainHandler := handlers.MainProgram{mockMain}
	mockPublisher := mainHandler.BuildPublisher()

	// Act
	controller := realMain.BuildController(mockPublisher)

	// Assert
	assert.NotNil(t, controller)
}

func TestRunSetsUpEnvironment(t *testing.T) {
	// Arrange
	mockMain := new(mocks.MainProgramInterface)
	realMain := new(handlers.MainProgram)
	mainHandler := handlers.MainProgram{mockMain}
	mockPublisher := mainHandler.BuildPublisher()
	mockController := mainHandler.BuildController(mockPublisher)
	mockConsumer := mainHandler.BuildConsumer(mockPublisher)

	// Act
	realMain.Run(mockPublisher, mockController, mockConsumer)
}
