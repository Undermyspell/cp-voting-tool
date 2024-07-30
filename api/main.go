package main

import (
	"net/http"
	"time"
	"voting/internal/env"
	"voting/shared/auth"
	"voting/shared/auth/jwks"
	"voting/shared/auth/middleware"
	shared_infra "voting/shared/infra/broker"
	user_http "voting/user/interface/http"
	votinghttp "voting/voting/interface/http"
	voting_sse "voting/voting/interface/sse"
	voting_ws "voting/voting/interface/ws"
	voting_repositories "voting/voting/repositories"

	_ "voting/docs"

	"github.com/centrifugal/centrifuge"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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
	centrifugeBroker := voting_ws.NewCentrifuge(internalBroker)

	r = gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	var votingStorage voting_repositories.VotingStorage
	if env.Env.VotingStorageInMemory {
		votingStorage = voting_repositories.NewInMemory()
	} else {
		logrus.Info("We use Postgresql")
		votingStorage = voting_repositories.NewPostgresql()
	}

	voting_repositories.InitInstances(votingStorage)

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
		v1.GET("/events", middleware.GinRequireAuth(), voting_sse.SseStream(internalBroker))
		q := v1.Group("/question", middleware.GinRequireAuth())
		q.PUT("/answer/:id", middleware.RequireRole(auth.SessionAdmin, auth.Admin), votinghttp.Answer)
		q.POST("/new", votinghttp.Create)
		q.PUT("/upvote/:id", votinghttp.Upvote)
		q.PUT("/undovote/:id", votinghttp.UndoVote)
		q.PUT("/update", votinghttp.Update)
		q.DELETE("/delete/:id", votinghttp.Delete)

		s := q.Group("/session", middleware.GinRequireAuth())
		s.POST("/start", middleware.RequireRole(auth.Admin), votinghttp.StartSession)
		s.POST("/stop", middleware.RequireRole(auth.Admin), votinghttp.StopSession)
		s.GET("", votinghttp.GetSession)

		ut := v1.Group("/user/test")
		ut.POST("/contributor", user_http.GetContributor)
		ut.POST("/admin", user_http.GetAdmin)
		ut.POST("/sessionadmin", user_http.GetAdmin)
	}

	start(r)
}
