package interfaces

//go:generate moq -out ./util_mocks/publisher_test.go . publisherInterface
type PublisherInterface interface {
	PublishMessage(exchangeName string, routingKey string, payload interface{}) error
	StopPublisher() error
	InitializePublisher() error
}
