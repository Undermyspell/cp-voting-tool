package main

import (
	"sse/internal/broker"
	"sse/internal/jwks"
	"sse/internal/middleware"
	services "sse/internal/services/question"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func addQuestionRoutes(r *gin.Engine, questionService services.QuestionService) {
	q := r.Group("/question")
	q.POST("/new", questionService.AddQuestion)
	q.PUT("/answer", questionService.Answer)
	q.PUT("/upvote", questionService.UpvoteQuestion)
	q.POST("/reset", questionService.Reset)

}

var initJwks = func() jwks.KeyfuncProvider {
	return jwks.New()
}
var start = func(r *gin.Engine) {
	r.Run(":3333")
}
var r *gin.Engine

func main() {
	r = gin.Default()

	jwks := initJwks()
	broker := broker.New()
	questionService := services.NewBrokered(broker)

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))
	r.Use(middleware.RequireAuth(jwks))

	r.GET("/events", broker.Stream)

	addQuestionRoutes(r, questionService)

	start(r)
}
