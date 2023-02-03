package broker

import "sse/internal/sse"

func New() Broker {
	return &ChannelBroker{
		Notifier:       make(chan sse.Event, 1),
		NewClients:     make(chan chan sse.Event),
		ClosingClients: make(chan chan sse.Event),
		Clients:        make(map[chan sse.Event]bool),
	}
}
