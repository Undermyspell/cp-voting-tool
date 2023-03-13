package broker

import (
	"fmt"
	"io"
	"sse/internal/models"
	"sse/internal/sse"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserBoundSseChannel struct {
	Channel chan sse.Event
	User    models.UserContext
}

type ChannelBroker struct {
	NotifierAll        chan sse.Event
	NotifierUser       chan sse.UserBoundSseEvent
	NotifierAllButUser chan sse.UserBoundSseEvent
	NewClients         chan UserBoundSseChannel
	ClosingClients     chan UserBoundSseChannel
	Clients            map[UserBoundSseChannel]bool
}

func (broker *ChannelBroker) Stream(c *gin.Context) {
	user, _ := c.Get(models.User)
	userContext := user.(*models.UserContext)
	userBoundSseChannel := UserBoundSseChannel{
		Channel: make(chan sse.Event),
		User:    *userContext,
	}

	defer func() {
		broker.ClosingClients <- userBoundSseChannel
		close(userBoundSseChannel.Channel)
	}()

	broker.NewClients <- userBoundSseChannel

	c.Stream(func(w io.Writer) bool {
		select {
		case event := <-userBoundSseChannel.Channel:
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
			logrus.Infof("🟢 Client added. %d registered clients", len(broker.Clients))

			event := sse.Event{
				EventType: sse.CLIENT_CONNECT_CHANGE,
				Payload:   fmt.Sprint(len(broker.Clients)),
			}

			broker.NotifyAll(event)
		case s := <-broker.ClosingClients:
			delete(broker.Clients, s)
			logrus.Infof("🔴 Removed client. %d registered clients", len(broker.Clients))
			event := sse.Event{
				EventType: sse.CLIENT_CONNECT_CHANGE,
				Payload:   fmt.Sprint(len(broker.Clients)),
			}

			broker.NotifyAll(event)
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

func (broker *ChannelBroker) NotifyAll(event sse.Event) {
	broker.NotifierAll <- event
}

func (broker *ChannelBroker) NotifyUser(event sse.Event, user models.UserContext) {
	broker.NotifierUser <- sse.UserBoundSseEvent{
		Event: event,
		User:  user,
	}
}

func (broker *ChannelBroker) NotifyAllButUser(event sse.Event, user models.UserContext) {
	broker.NotifierAllButUser <- sse.UserBoundSseEvent{
		Event: event,
		User:  user,
	}
}
