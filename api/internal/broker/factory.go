package broker

import "voting/internal/events"

func New() Broker {

	broker := &ChannelBroker{
		NotifierAll:        make(chan events.Event, 1),
		NotifierUser:       make(chan events.UserBoundSseEvent, 1),
		NotifierAllButUser: make(chan events.UserBoundSseEvent, 1),
		NewClients:         make(chan UserBoundSseChannel),
		ClosingClients:     make(chan UserBoundSseChannel),
		Clients:            make(map[UserBoundSseChannel]bool),
	}

	go broker.Listen()
	go broker.SendHeartBeat()

	return broker
}
