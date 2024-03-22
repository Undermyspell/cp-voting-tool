package main

import (
	"net/http"
	"time"
	"voting/internal/env"
	"voting/internal/jwks"
	"voting/internal/middleware"
	"voting/internal/models/roles"
	"voting/internal/notification"
	"voting/internal/votingstorage"
	questionService "voting/services/question"
	userService "voting/services/user"
	shared_infra "voting/shared/infra/broker"
	votinghttp "voting/voting/interface/http"
	voting_repositories "voting/voting/repositories"

	_ "voting/docs"

	"github.com/centrifugal/centrifuge"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var initQuestionService = func(broker shared_infra.Broker, votingStorage votingstorage.VotingStorage) questionService.QuestionService {
	return questionService.NewBrokered(broker, votingStorage)
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
	jwks.Init()

	internalBroker := shared_infra.New()
	centrifugeBroker := notification.NewCentrifuge(internalBroker)

	r = gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	var votingStorage votingstorage.VotingStorage
	if env.Env.VotingStorageInMemory {
		votingStorage = votingstorage.NewInMemory()
	} else {
		logrus.Info("WE USE REDIS")
		votingStorage = votingstorage.NewRedis()
	}

	voting_repositories.InitInstances(votingStorage)

	questionService := initQuestionService(internalBroker, votingStorage)
	userService := userService.NewTestUser()

	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))
	r.Use(middleware.Options)
	r.Use(middleware.GinContextToContextMiddleware())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/connection/websocket", gin.WrapH(middleware.CentrifugeAnonymousAuth(centrifuge.NewWebsocketHandler(centrifugeBroker, centrifuge.WebsocketConfig{
			CheckOrigin: func(r *http.Request) bool {
				originHeader := r.Header.Get("Origin")
				if originHeader == "" {
					return true
				}
				return originHeader == env.Env.AllowedOrigin
			},
		}))))
		v1.GET("/events", middleware.GinRequireAuth(), notification.SseStream(internalBroker))
		q := v1.Group("/question", middleware.GinRequireAuth())
		q.PUT("/answer/:id", middleware.RequireRole(roles.SessionAdmin, roles.Admin), questionService.Answer)
		q.POST("/new", votinghttp.Create)
		q.PUT("/upvote/:id", votinghttp.Upvote)
		q.PUT("/undovote/:id", questionService.UndoVote)
		q.PUT("/update", votinghttp.Update)
		q.DELETE("/delete/:id", questionService.Delete)

		s := q.Group("/session", middleware.GinRequireAuth())
		s.POST("/start", middleware.RequireRole(roles.Admin), votinghttp.StartSession)
		s.POST("/stop", middleware.RequireRole(roles.Admin), votinghttp.StopSession)
		s.GET("", questionService.GetSession)

		ut := v1.Group("/user/test")
		ut.POST("/contributor", userService.GetContributor)
		ut.POST("/admin", userService.GetAdmin)
		ut.POST("/sessionadmin", userService.GetAdmin)
	}

	start(r)
}
