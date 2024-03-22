package shared_infra_broker

import (
	"voting/internal/events"
	"voting/shared/shared_models"
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
	NotifyAll(events.Event)
	NotifyUser(event events.Event, user shared_models.UserContext)
	NotifyAllButUser(event events.Event, user shared_models.UserContext)
	DistinctClientsCount() int
	SendHeartBeat()
	AddClient(client UserBoundChannel)
	RemoveClient(client UserBoundChannel)
}
