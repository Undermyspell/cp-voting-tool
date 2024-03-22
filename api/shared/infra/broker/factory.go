package shared_infra_broker

import (
	"voting/internal/events"
)

func New() Broker {

	broker := &InternalBroker{
		NotifierAll:        make(chan events.Event, 1),
		NotifierUser:       make(chan events.UserBoundEvent, 1),
		NotifierAllButUser: make(chan events.UserBoundEvent, 1),
		NewClients:         make(chan UserBoundChannel),
		ClosingClients:     make(chan UserBoundChannel),
		Clients:            make(map[UserBoundChannel]bool),
	}

	go broker.Listen()
	initInstances(broker)

	return broker
}
