package main

import (
	"sse/internal/broker"
	"sse/internal/env"
	"sse/internal/jwks"
	"sse/internal/middleware"
	"sse/internal/mocks"
	"sse/internal/models/roles"
	questionService "sse/services/question"
	userService "sse/services/user"
	"time"

	_ "sse/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var initQuestionService = func(broker broker.Broker) questionService.QuestionService {
	return questionService.NewBrokered(broker)
}

var start = func(r *gin.Engine) {
	r.Run(":3333")
}
var r *gin.Engine

// @title           CP Voting tool api
// @version         1.0
// @description     THE CP voting tool API in Go using Gin framework.

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
// @host      localhost:3333
// @BasePath  /api/v1
func main() {
	env.Init()

	r = gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	var jwksProvider jwks.KeyfuncProvider
	if !env.Env.UseMockJwks {
		jwksProvider = jwks.New()
	} else {
		jwksProvider = mocks.NewJwks()
	}

	broker := broker.New()
	questionService := initQuestionService(broker)
	userService := userService.NewTestUser()

	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/events", middleware.RequireAuth(jwksProvider), broker.Stream)
		q := v1.Group("/question", middleware.RequireAuth(jwksProvider))
		q.PUT("/answer/:id", middleware.RequireRole(roles.SessionAdmin, roles.Admin), questionService.Answer)
		q.POST("/new", questionService.Add)
		q.PUT("/upvote/:id", questionService.Upvote)
		q.PUT("/update", questionService.Update)
		q.DELETE("/delete/:id", questionService.Delete)

		s := q.Group("/session", middleware.RequireAuth(jwksProvider))
		s.POST("/start", middleware.RequireRole(roles.Admin), questionService.Start)
		s.POST("/stop", middleware.RequireRole(roles.Admin), questionService.Stop)
		s.GET("", questionService.GetSession)

		ut := v1.Group("/user/test")
		ut.POST("/contributor", userService.GetContributor)
		ut.POST("/admin", userService.GetAdmin)
		ut.POST("/sessionadmin", userService.GetAdmin)
	}

	start(r)
}
