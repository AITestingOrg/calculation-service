package interfaces

type PublisherInterface interface {
	PublishMessage(exchangeName string, routingKey string, payload interface{}) error
	StopPublisher() error
	InitializePublisher() error
}
