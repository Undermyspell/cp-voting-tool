package broker

import (
	"voting/internal/events"
	"voting/internal/models"

	"github.com/gin-gonic/gin"
)

type Broker interface {
	SseStream(c *gin.Context)
	Listen()
	NotifyAll(events.Event)
	NotifyUser(event events.Event, user models.UserContext)
	NotifyAllButUser(event events.Event, user models.UserContext)
	DistinctClientsCount() int
	SendHeartBeat()
}
