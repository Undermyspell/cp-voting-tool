package broker

import (
	"sse/internal/models"
	"sse/internal/sse"

	"github.com/gin-gonic/gin"
)

type Broker interface {
	Stream(c *gin.Context)
	Listen()
	NotifyAll(sse.Event)
	NotifyUser(event sse.Event, user models.UserContext)
	NotifyAllButUser(event sse.Event, user models.UserContext)
	DistinctClientsCount() int
	SendHeartBeat()
}
