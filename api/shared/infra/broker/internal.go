package shared_infra_broker

import (
	"encoding/json"
	"time"
	"voting/internal/events"
	"voting/internal/models"

	"github.com/sirupsen/logrus"
)

type UserBoundChannel struct {
	Channel chan events.Event
	User    models.UserContext
}

type InternalBroker struct {
	NotifierAll        chan events.Event
	NotifierUser       chan events.UserBoundEvent
	NotifierAllButUser chan events.UserBoundEvent
	NewClients         chan UserBoundChannel
	ClosingClients     chan UserBoundChannel
	Clients            map[UserBoundChannel]bool
}

func (broker *InternalBroker) Listen() {
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

func (broker *InternalBroker) NotifyAll(event events.Event) {
	broker.NotifierAll <- event
}

func (broker *InternalBroker) NotifyUser(event events.Event, user models.UserContext) {
	broker.NotifierUser <- events.UserBoundEvent{
		Event: event,
		User:  user,
	}
}

func (broker *InternalBroker) NotifyAllButUser(event events.Event, user models.UserContext) {
	broker.NotifierAllButUser <- events.UserBoundEvent{
		Event: event,
		User:  user,
	}
}

func (broker *InternalBroker) DistinctClientsCount() int {
	distinctClients := make(map[models.UserContext]bool)
	for clientMessageChan := range broker.Clients {
		if _, ok := distinctClients[clientMessageChan.User]; !ok {
			distinctClients[clientMessageChan.User] = true
		}
	}

	return len(distinctClients)
}

func (broker *InternalBroker) SendHeartBeat() {
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

func (broker *InternalBroker) AddClient(client UserBoundChannel) {
	broker.NewClients <- client
}

func (broker *InternalBroker) RemoveClient(client UserBoundChannel) {
	broker.ClosingClients <- client
}

func (broker *InternalBroker) createUserConnectionEvent(eventType events.EventType) events.Event {
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
