package main

import (
	"sse/internal/broker"
	"sse/internal/jwks"
	"sse/internal/middleware"
	"sse/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	teardown := make(chan string)
	r := gin.Default()
	jwks.Init(teardown)

	broker := broker.New()
	go broker.Listen()

	questionService := services.New(broker)

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}

	r.Use(cors.New(config))
	r.Use(middleware.RequireAuth())

	r.GET("/events", broker.Stream)
	q := r.Group("/question")
	q.POST("/new", questionService.AddQuestion)
	q.PUT("/answer", questionService.Answer)
	q.PUT("/upvote", questionService.UpvoteQuestion)

	// r.POST("/reset", questionService.Reset)
	r.Run(":3333")

	teardown <- "teardown"
}
