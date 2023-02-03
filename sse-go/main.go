package main

import (
	"sse/internal/broker"
	"sse/internal/jwks"
	"sse/internal/middleware"
	services "sse/internal/services/question"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var initJwks = func() jwks.KeyfuncProvider {
	return jwks.New()
}

var initQuestionService = func(broker broker.Broker) services.QuestionService {
	return services.NewBrokered(broker)
}

var start = func(r *gin.Engine) {
	r.Run(":3333")
}
var r *gin.Engine

func main() {
	r = gin.Default()

	jwks := initJwks()
	broker := broker.New()
	questionService := initQuestionService(broker)

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))
	r.Use(middleware.RequireAuth(jwks))

	r.GET("/events", broker.Stream)

	q := r.Group("/question")
	q.POST("/new", questionService.AddQuestion)
	q.PUT("/answer", questionService.Answer)
	q.PUT("/upvote", questionService.UpvoteQuestion)
	q.POST("/reset", questionService.Reset)

	start(r)
}
