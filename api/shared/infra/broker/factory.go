package shared_infra_broker

import (
	usecases_events "voting/voting/use-cases/_events"
)

func New() Broker {

	broker := &InternalBroker{
		NotifierAll:        make(chan usecases_events.Event, 1),
		NotifierUser:       make(chan usecases_events.UserBoundEvent, 1),
		NotifierAllButUser: make(chan usecases_events.UserBoundEvent, 1),
		NewClients:         make(chan UserBoundChannel),
		ClosingClients:     make(chan UserBoundChannel),
		Clients:            make(map[UserBoundChannel]bool),
	}

	go broker.Listen()
	initInstances(broker)

	return broker
}
