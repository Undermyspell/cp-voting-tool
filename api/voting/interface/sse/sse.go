package voting_sse

import (
	"io"
	shared_infra_broker "voting/shared/infra/broker"
	"voting/shared/shared_models"
	usecases_events "voting/voting/use-cases/_events"

	"github.com/gin-gonic/gin"
)

func SseStream(internalBroker shared_infra_broker.Broker) func(c *gin.Context) {
	return func(c *gin.Context) {
		user, _ := c.Get(shared_models.User)
		userContext := user.(*shared_models.UserContext)

		userBoundChannel := shared_infra_broker.UserBoundChannel{
			Channel: make(chan usecases_events.Event),
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
