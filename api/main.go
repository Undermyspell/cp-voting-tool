package main

import (
	"net/http"
	"sse/internal/broker"
	"sse/internal/env"
	"sse/internal/jwks"
	"sse/internal/middleware"
	"sse/internal/mocks"
	services "sse/services/question"
	"time"

	_ "sse/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var initQuestionService = func(broker broker.Broker) services.QuestionService {
	return services.NewBrokered(broker)
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

	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))
	// r.Use(middleware.RequireAuth(jwksProvider))

	v1 := r.Group("/api/v1", middleware.RequireAuth(jwksProvider))
	{
		v1.GET("/events", broker.Stream)
		q := v1.Group("/question")
		q.PUT("/answer/:id", questionService.Answer)
		q.POST("/new", questionService.Add)
		q.PUT("/upvote/:id", questionService.Upvote)
		q.PUT("/update", questionService.Update)
		q.DELETE("/delete/:id", questionService.Delete)

		s := q.Group("/session")
		s.POST("/start", questionService.Start)
		s.POST("/stop", questionService.Stop)
		s.GET("", questionService.GetSession)
	}
	r.POST("mockuser", func(c *gin.Context) {
		usr := struct {
			FirstName string
			LastName  string
		}{}

		err := c.BindJSON(&usr)

		if err != nil {
			logrus.Error(err.Error())
			c.JSON(http.StatusBadRequest, "cant parse request")
			return
		}

		token := mocks.GetToken(usr.FirstName, usr.LastName)

		c.JSON(http.StatusOK, struct{ Token string }{Token: token})
	})

	start(r)
}
