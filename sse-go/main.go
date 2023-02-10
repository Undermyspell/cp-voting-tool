package main

import (
	"sse/internal/broker"
	"sse/internal/jwks"
	"sse/internal/middleware"
	"sse/internal/models/roles"
	services "sse/services/question"

	_ "sse/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

// @title           Voting tool api
// @version         1.0
// @description     A voting tool API in Go using Gin framework.

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
// @host      localhost:3333
// @BasePath  /api/v1
func main() {
	r = gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	jwks := initJwks()
	broker := broker.New()
	questionService := initQuestionService(broker)

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))
	r.Use(middleware.RequireAuth(jwks))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/events", broker.Stream)
		q := v1.Group("/question")
		q.PUT("/answer/:id", middleware.RequireRole(roles.SessionAdmin), questionService.Answer)
		q.POST("/new", questionService.AddQuestion)
		q.PUT("/upvote/:id", questionService.UpvoteQuestion)
		q.POST("/reset", questionService.Reset)
		q.GET("/session", questionService.GetSession)
	}

	start(r)
}
