package shared_infra_broker

import (
	"voting/shared"
	shared_models "voting/shared/models"
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
	NotifyAll(shared.Event)
	NotifyUser(event shared.Event, user shared_models.UserContext)
	NotifyAllButUser(event shared.Event, user shared_models.UserContext)
	DistinctClientsCount() int
	SendHeartBeat()
	AddClient(client UserBoundChannel)
	RemoveClient(client UserBoundChannel)
}
