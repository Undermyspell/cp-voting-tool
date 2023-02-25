package broker

import (
	"io"
	"sse/internal/sse"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ChannelBroker struct {
	Notifier       chan sse.Event
	NewClients     chan chan sse.Event
	ClosingClients chan chan sse.Event
	Clients        map[chan sse.Event]bool
}

func (broker *ChannelBroker) Stream(c *gin.Context) {
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

func (broker *ChannelBroker) Listen() {
	for {
		select {
		case s := <-broker.NewClients:
			// A new client has joined
			broker.Clients[s] = true
			logrus.Infof("🟢 Client added. %d registered clients", len(broker.Clients))
		case s := <-broker.ClosingClients:
			// A client has detached
			// remove them from our clients map
			delete(broker.Clients, s)
			logrus.Infof("🔴 Removed client. %d registered clients", len(broker.Clients))
		case event := <-broker.Notifier:
			// case for getting a new msg
			// Thus send it to all clients
			for clientMessageChan := range broker.Clients {
				clientMessageChan <- event
			}
		}
	}
}

func (broker *ChannelBroker) Notify(event sse.Event) {
	broker.Notifier <- event
}