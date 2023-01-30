package main

import (
	"sse/internal/broker"
	"sse/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	broker := broker.New()
	go broker.Listen()

	questionService := services.New(broker)

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}

	r.Use(cors.New(config))

	r.GET("/events", broker.Stream)
	q := r.Group("/question")
	q.POST("/new", questionService.AddQuestion)
	q.PUT("/answer", questionService.Answer)
	q.PUT("/upvote", questionService.UpvoteQuestion)
	r.POST("/reset", questionService.Reset)
	r.Run(":3333")
}
