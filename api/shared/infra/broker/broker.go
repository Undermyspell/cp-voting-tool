package shared_infra_broker

import (
	shared_models "voting/shared/models"
	usecases_events "voting/voting/use-cases/_events"
)

var b Broker

func GetInstance() Broker {
	return b
}

func initInstances(broker Broker) {
	b = broker
}

type Broker interface {
	Listen()
	NotifyAll(usecases_events.Event)
	NotifyUser(event usecases_events.Event, user shared_models.UserContext)
	NotifyAllButUser(event usecases_events.Event, user shared_models.UserContext)
	DistinctClientsCount() int
	SendHeartBeat()
	AddClient(client UserBoundChannel)
	RemoveClient(client UserBoundChannel)
}
