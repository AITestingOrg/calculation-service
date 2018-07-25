package interfaces

type MainProgramInterface interface {
	BuildPublisher() PublisherInterface
	BuildConsumer(publisher *PublisherInterface) ConsumerInterface
	BuildController(publisher *PublisherInterface) ControllerInterface
	Run(publisher PublisherInterface, controller ControllerInterface, consumer ConsumerInterface)
}
