package shared_infra_broker

import (
	"voting/shared"
)

func New() Broker {

	broker := &InternalBroker{
		NotifierAll:        make(chan shared.Event, 1),
		NotifierUser:       make(chan shared.UserBoundEvent, 1),
		NotifierAllButUser: make(chan shared.UserBoundEvent, 1),
		NewClients:         make(chan UserBoundChannel),
		ClosingClients:     make(chan UserBoundChannel),
		Clients:            make(map[UserBoundChannel]bool),
	}

	go broker.Listen()
	initInstances(broker)

	return broker
}
