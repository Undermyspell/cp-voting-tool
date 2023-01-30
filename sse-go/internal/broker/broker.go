package broker

import (
	"io"
	"log"
	"sse/internal/sse"

	"github.com/gin-gonic/gin"
)

type Broker struct {
	Notifier       chan sse.Event
	NewClients     chan chan sse.Event
	ClosingClients chan chan sse.Event
	Clients        map[chan sse.Event]bool
}

func New() *Broker {
	return &Broker{
		Notifier:       make(chan sse.Event, 1),
		NewClients:     make(chan chan sse.Event),
		ClosingClients: make(chan chan sse.Event),
		Clients:        make(map[chan sse.Event]bool),
	}
}

func (broker *Broker) Stream(c *gin.Context) {
	messageChan := make(chan sse.Event)

	defer func() {
		broker.ClosingClients <- messageChan
		close(messageChan)
	}()

	broker.NewClients <- messageChan

	c.Stream(func(w io.Writer) bool {
		select {
		case event := <-messageChan:
			c.SSEvent(string(event.EventType), event.Payload)
		case <-c.Request.Context().Done():
			return false
		}
		return true
	})
}

func (broker *Broker) Listen() {
	for {
		select {
		case s := <-broker.NewClients:
			// A new client has joined
			broker.Clients[s] = true
			log.Printf("🟢 Client added. %d registered clients", len(broker.Clients))
		case s := <-broker.ClosingClients:
			// A client has detached
			// remove them from our clients map
			delete(broker.Clients, s)
			log.Printf("🔴 Removed client. %d registered clients", len(broker.Clients))
		case event := <-broker.Notifier:
			// case for getting a new msg
			// Thus send it to all clients
			for clientMessageChan := range broker.Clients {
				clientMessageChan <- event
			}
		}
	}
}
