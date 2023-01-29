package main

import (
	"sse/internal/broker"
	"sse/internal/events"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Widget struct {
	Id    int
	Attrs []string
}

func main() {
	r := gin.Default()

	broker := broker.New()
	go broker.Listen()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}

	r.Use(cors.New(config))

	// r.GET("/events", broker.Stream)
	r.GET("/events", broker.StreamNew(events.NEW_QUESTION))
	r.POST("/send", broker.BroadcastMessage)
	r.Run(":3333")
}
