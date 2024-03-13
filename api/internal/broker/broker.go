package broker

import (
	"voting/internal/events"
	"voting/internal/models"
)

type Broker interface {
	Listen()
	NotifyAll(events.Event)
	NotifyUser(event events.Event, user models.UserContext)
	NotifyAllButUser(event events.Event, user models.UserContext)
	DistinctClientsCount() int
	SendHeartBeat()
	AddClient(client UserBoundChannel)
	RemoveClient(client UserBoundChannel)
}
