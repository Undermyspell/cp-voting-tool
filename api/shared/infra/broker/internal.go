package shared_infra_broker

import (
	"encoding/json"
	"time"
	"voting/shared"
	shared_models "voting/shared/models"
	usecases_events "voting/voting/use-cases/_events"

	"github.com/sirupsen/logrus"
)

type UserBoundChannel struct {
	Channel chan shared.Event
	User    shared_models.UserContext
}

type InternalBroker struct {
	NotifierAll        chan shared.Event
	NotifierUser       chan shared.UserBoundEvent
	NotifierAllButUser chan shared.UserBoundEvent
	NewClients         chan UserBoundChannel
	ClosingClients     chan UserBoundChannel
	Clients            map[UserBoundChannel]bool
}

func (broker *InternalBroker) Listen() {
	for {
		select {
		case s := <-broker.NewClients:
			broker.Clients[s] = true
			event := broker.createUserConnectionEvent(usecases_events.USER_CONNECTED)
			for clientMessageChan := range broker.Clients {
				clientMessageChan.Channel <- event
			}
			logrus.Infof("🟢 Client added. %d registered clients", len(broker.Clients))
		case s := <-broker.ClosingClients:
			delete(broker.Clients, s)
			event := broker.createUserConnectionEvent(usecases_events.USER_DISCONNECTED)
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

func (broker *InternalBroker) NotifyAll(event shared.Event) {
	broker.NotifierAll <- event
}

func (broker *InternalBroker) NotifyUser(event shared.Event, user shared_models.UserContext) {
	broker.NotifierUser <- shared.UserBoundEvent{
		Event: event,
		User:  user,
	}
}

func (broker *InternalBroker) NotifyAllButUser(event shared.Event, user shared_models.UserContext) {
	broker.NotifierAllButUser <- shared.UserBoundEvent{
		Event: event,
		User:  user,
	}
}

func (broker *InternalBroker) DistinctClientsCount() int {
	distinctClients := make(map[shared_models.UserContext]bool)
	for clientMessageChan := range broker.Clients {
		if _, ok := distinctClients[clientMessageChan.User]; !ok {
			distinctClients[clientMessageChan.User] = true
		}
	}

	return len(distinctClients)
}

func (broker *InternalBroker) SendHeartBeat() {
	event := shared.Event{
		EventType: usecases_events.HEART_BEAT,
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

func (broker *InternalBroker) createUserConnectionEvent(eventType shared.EventType) shared.Event {
	event := usecases_events.UserConnected{
		UserCount: broker.DistinctClientsCount(),
	}
	eventByteArray, _ := json.Marshal(event)
	disconnectedEvent := shared.Event{
		EventType: eventType,
		Payload:   string(eventByteArray),
	}

	return disconnectedEvent
}
