package main

import (
	"context"
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"time"
	"voting/internal/env"
	"voting/shared/auth"
	authhandler "voting/shared/auth/handler"
	"voting/shared/auth/jwks"
	"voting/shared/auth/middleware"
	authsession "voting/shared/auth/session"
	shared_infra "voting/shared/infra/broker"
	"voting/templates"
	"voting/templates/components"
	user_http "voting/user/interface/http"
	votinghttp "voting/voting/interface/http"
	voting_sse "voting/voting/interface/sse"
	voting_ws "voting/voting/interface/ws"
	voting_repositories "voting/voting/repositories"
	voting_usecases "voting/voting/usecases"

	_ "voting/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed static/*
var static embed.FS

//go:generate npx tailwindcss build -i static/css/style.css -o static/css/tailwind.css -m

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

	store := authsession.InitSessionStore()
	r.Use(sessions.Sessions("auth-session", store))

	fsys, _ := fs.Sub(static, "static")
	r.StaticFS("/static", http.FS(fsys))
	r.NoRoute(NoRouteHanlder)

	app := r.Group("/", middleware.GinRequireCookieAuth())
	{
		app.GET("/", homeHandler)
		app.GET("/login", authhandler.Login)
		app.GET("/oauth2/callback", authhandler.LoginCallback)
		app.GET("/user", user_http.GetAuthenticatedUser)
		app.GET("/q/new", newQuestionHandler)
		app.POST("/q/save", func(c *gin.Context) {
			c.Header("HX-Redirect", "/")
			c.Status(http.StatusCreated)
		})
	}

	api := r.Group("/api/v1")
	{
		api.GET("/connection/websocket", votinghttp.CentrifugoHandler())
		api.GET("/events", middleware.GinRequireJwtAuth(), voting_sse.SseStream(internalBroker))
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

func NoRouteHanlder(c *gin.Context) {
	logrus.Info("NO ROUTE")
	c.Redirect(http.StatusMovedPermanently, "/")
}

func newQuestionHandler(c *gin.Context) {
	component := components.NewQuestionModal()
	component.Render(c.Request.Context(), c.Writer)
}

func homeHandler(c *gin.Context) {
	sessions := sessions.Default(c)
	token := sessions.Get("token").(string)

	req, err := http.NewRequest("GET", "http://:3333/api/v1/question/session", nil)

	if err != nil {
		logrus.Error(err)
		c.String(http.StatusInternalServerError, "Error fetching questions")
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	// Create an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		c.String(http.StatusInternalServerError, "Error fetching questions")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.String(http.StatusInternalServerError, "Error fetching questions")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		c.String(http.StatusInternalServerError, "Error fetching questions")
		return
	}

	var questions []voting_usecases.QuestionDto
	err = json.Unmarshal(body, &questions)
	if err != nil {
		logrus.Error(err)
		c.String(http.StatusInternalServerError, "Error fetching questions")
	}

	component := templates.Main("Home Page :)", "Welcome to the Home Page :)!", questions)
	component.Render(c.Request.Context(), c.Writer)
}
