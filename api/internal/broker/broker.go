package broker

import (
	"sse/internal/sse"

	"github.com/gin-gonic/gin"
)

type Broker interface {
	Stream(c *gin.Context)
	Listen()
	Notify(sse.Event)
}
