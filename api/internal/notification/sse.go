package notification

import (
	"io"
	"voting/internal/broker"
	"voting/internal/events"
	"voting/internal/models"

	"github.com/gin-gonic/gin"
)

func SseStream(internalBroker broker.Broker) func(c *gin.Context) {
	return func(c *gin.Context) {
		user, _ := c.Get(models.User)
		userContext := user.(*models.UserContext)

		userBoundChannel := broker.UserBoundChannel{
			Channel: make(chan events.Event),
			User:    *userContext,
		}

		defer func() {
			internalBroker.RemoveClient(userBoundChannel)
			close(userBoundChannel.Channel)
		}()

		internalBroker.AddClient(userBoundChannel)

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
}
