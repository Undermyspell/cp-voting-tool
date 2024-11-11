package main

import (
	"embed"
	"io/fs"
	"net/http"
	bff "voting/bff/interface/http"
	"voting/internal/env"
	"voting/shared/auth"
	authhandler "voting/shared/auth/handler"
	"voting/shared/auth/jwks"
	"voting/shared/auth/middleware"
	authsession "voting/shared/auth/session"
	shared_infra "voting/shared/infra/broker"
	user_http "voting/user/interface/http"
	votinghttp "voting/voting/interface/http"
	voting_ws "voting/voting/interface/ws"
	voting_repositories "voting/voting/repositories"

	_ "voting/docs"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed bff/static/*
var static embed.FS

//go:generate npx tailwindcss build -i bff/static/css/style.css -o bff/static/css/tailwind.css -m
//go:generate npx esbuild --bundle --outfile=bff/static/js/index.js bff/client/index.ts

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
	authhandler.InitOAuthConfig()

	internalBroker := shared_infra.New()
	voting_ws.InitCentrifuge(internalBroker)

	r = gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	var votingStorage voting_repositories.VotingStorage

	switch env.Env.Storage {
	case env.InMemory:
		votingStorage = voting_repositories.NewInMemory()
	case env.Postgres:
		votingStorage = voting_repositories.NewPostgresql()
	case env.Redis:
		votingStorage = voting_repositories.NewRedis()
	}

	logrus.Infof("Using storage: %s", env.Env.Storage)

	voting_repositories.InitInstances(votingStorage)

	r.Use(middleware.Cors())
	r.Use(middleware.Options)
	r.Use(middleware.GinContextToContextMiddleware())

	store := authsession.InitSessionStore()
	r.Use(sessions.Sessions("auth-session", store))

	fsys, _ := fs.Sub(static, "bff/static")
	r.StaticFS("/static", http.FS(fsys))
	r.NoRoute(NoRouteHandler)
	app := r.Group("/", middleware.GinRequireCookieAuth())
	{
		app.GET("/", bff.Main)
		app.GET("/login", authhandler.Login)
		app.GET("/oauth2/callback", authhandler.LoginCallback)
		app.GET("/user", user_http.GetAuthenticatedUser)
		app.GET("/q/new", bff.NewQuestion)
		app.GET("/q/update/:id", bff.UpdateQuestion)
		app.POST("/q/save", bff.SaveNewQuestion)
		app.PUT("/q/save", bff.SaveUpdatedQuestion)
		app.DELETE("/q/delete/:id", bff.DeleteQuestion)
		app.PUT("/q/upvote/:id", bff.UpvoteQuestion)
		app.PUT("/q/undovote/:id", bff.UndoVoteQuestion)
		app.PUT("/q/answer/:id", bff.AnswerQuestion)
		app.POST("/q/s/start", bff.StartSession)
		app.POST("/q/s/stop", bff.StopSession)
		app.GET("/q/s/page/:activeSession/:onlyUnanswered", bff.QuestionSessionPage)
		app.GET("/q/s/download", bff.DownloadSessionAsCsv)
	}

	api := r.Group("/api/v1")
	{
		api.GET("/connection/websocket", votinghttp.CentrifugoHandler())
		q := api.Group("/question", middleware.GinRequireJwtAuth())
		q.PUT("/answer/:id", middleware.RequireRole(auth.SessionAdmin, auth.Admin), votinghttp.Answer)
		q.POST("/new", votinghttp.Create)
		q.PUT("/upvote/:id", votinghttp.Upvote)
		q.PUT("/undovote/:id", votinghttp.UndoVote)
		q.PUT("/update", votinghttp.Update)
		q.DELETE("/delete/:id", votinghttp.Delete)

		s := q.Group("/session", middleware.GinRequireJwtAuth())
		s.POST("/start", middleware.RequireRole(auth.Admin), votinghttp.StartSession)
		s.POST("/stop", middleware.RequireRole(auth.Admin), votinghttp.StopSession)
		s.GET("", votinghttp.GetSession)

		ut := api.Group("/user/test")
		ut.POST("/contributor", user_http.GetContributor)
		ut.POST("/admin", user_http.GetAdmin)
		ut.POST("/sessionadmin", user_http.GetAdmin)
	}

	start(r)
}

func NoRouteHandler(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/")
}
