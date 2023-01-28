package main

import (
	"sse/broker"

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

	r.GET("/events", broker.Stream)
	r.POST("/send", broker.BroadcastMessage)
	r.Run(":3333")
}
