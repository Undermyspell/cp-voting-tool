package broker

import (
	"encoding/json"
	"io"
	"time"
	"voting/internal/events"
	"voting/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserBoundChannel struct {
	Channel chan events.Event
	User    models.UserContext
}

type ChannelBroker struct {
	NotifierAll        chan events.Event
	NotifierUser       chan events.UserBoundEvent
	NotifierAllButUser chan events.UserBoundEvent
	NewClients         chan UserBoundChannel
	ClosingClients     chan UserBoundChannel
	Clients            map[UserBoundChannel]bool
}

func (broker *ChannelBroker) SseStream(c *gin.Context) {
	user, _ := c.Get(models.User)
	userContext := user.(*models.UserContext)
	userBoundChannel := UserBoundChannel{
		Channel: make(chan events.Event),
		User:    *userContext,
	}

	defer func() {
		broker.ClosingClients <- userBoundChannel
		close(userBoundChannel.Channel)
	}()

	broker.NewClients <- userBoundChannel

	c.Stream(func(w io.Writer) bool {
		select {
		case event := <-userBoundChannel.Channel:
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
			broker.Clients[s] = true
			event := broker.createUserConnectionEvent(events.USER_CONNECTED)
			for clientMessageChan := range broker.Clients {
				clientMessageChan.Channel <- event
			}
			logrus.Infof("🟢 Client added. %d registered clients", len(broker.Clients))
		case s := <-broker.ClosingClients:
			delete(broker.Clients, s)
			event := broker.createUserConnectionEvent(events.USER_DISCONNECTED)
			for clientMessageChan := range broker.Clients {
				clientMessageChan.Channel <- event
			}
			logrus.Infof("🔴 Removed client. %d registered clients", len(broker.Clients))
		case event := <-broker.NotifierAll:
			for clientMessageChan := range broker.Clients {
				clientMessageChan.Channel <- event
			}
		case event := <-broker.NotifierUser:
			for clientMessageChan := range broker.Clients {
				if clientMessageChan.User.Email == event.User.Email {
					clientMessageChan.Channel <- event.Event
				}
			}
		case event := <-broker.NotifierAllButUser:
			for clientMessageChan := range broker.Clients {
				if clientMessageChan.User.Email != event.User.Email {
					clientMessageChan.Channel <- event.Event
				}
			}
		}
	}
}

func (broker *ChannelBroker) NotifyAll(event events.Event) {
	broker.NotifierAll <- event
}

func (broker *ChannelBroker) NotifyUser(event events.Event, user models.UserContext) {
	broker.NotifierUser <- events.UserBoundEvent{
		Event: event,
		User:  user,
	}
}

func (broker *ChannelBroker) NotifyAllButUser(event events.Event, user models.UserContext) {
	broker.NotifierAllButUser <- events.UserBoundEvent{
		Event: event,
		User:  user,
	}
}

func (broker *ChannelBroker) DistinctClientsCount() int {
	distinctClients := make(map[models.UserContext]bool)
	for clientMessageChan := range broker.Clients {
		if _, ok := distinctClients[clientMessageChan.User]; !ok {
			distinctClients[clientMessageChan.User] = true
		}
	}

	return len(distinctClients)
}

func (broker *ChannelBroker) SendHeartBeat() {
	event := events.Event{
		EventType: events.HEART_BEAT,
		Payload:   "",
	}

	ticker := time.NewTicker(30 * 1000 * time.Millisecond)
	for range ticker.C {
		logrus.Info("Send heartbeat...")
		broker.NotifyAll(event)
	}
}

func (broker *ChannelBroker) createUserConnectionEvent(eventType events.EventType) events.Event {
	event := events.UserConnected{
		UserCount: broker.DistinctClientsCount(),
	}
	eventByteArray, _ := json.Marshal(event)
	disconnectedEvent := events.Event{
		EventType: eventType,
		Payload:   string(eventByteArray),
	}

	return disconnectedEvent
}
