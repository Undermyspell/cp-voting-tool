package broker

import "sse/internal/sse"

func New() Broker {

	broker := &ChannelBroker{
		NotifierAll:        make(chan sse.Event, 1),
		NotifierUser:       make(chan sse.UserBoundSseEvent, 1),
		NotifierAllButUser: make(chan sse.UserBoundSseEvent, 1),
		NewClients:         make(chan UserBoundSseChannel),
		ClosingClients:     make(chan UserBoundSseChannel),
		Clients:            make(map[UserBoundSseChannel]bool),
	}

	go broker.Listen()

	return broker
}
