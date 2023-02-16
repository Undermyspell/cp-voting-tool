package main

import (
	"net/http"
	"os"
	"sse/internal/broker"
	"sse/internal/env"
	"sse/internal/jwks"
	"sse/internal/middleware"
	"sse/internal/mocks"
	"sse/internal/models/roles"
	services "sse/services/question"

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
	logrus.Printf("env %s", os.Getenv("APP_ENV"))
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

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))
	// r.Use(middleware.RequireAuth(jwksProvider))

	v1 := r.Group("/api/v1", middleware.RequireAuth(jwksProvider))
	{
		v1.GET("/events", broker.Stream)
		q := v1.Group("/question")
		q.PUT("/answer/:id", middleware.RequireRole(roles.SessionAdmin), questionService.Answer)
		q.POST("/new", questionService.AddQuestion)
		q.PUT("/upvote/:id", questionService.UpvoteQuestion)
		q.POST("/reset", questionService.Reset)
		q.GET("/session", questionService.GetSession)
	}
	r.GET("mockuser", func(c *gin.Context) {
		usr := struct {
			firstName string
			lastName  string
		}{}

		err := c.BindJSON(&usr)

		if err != nil {
			logrus.Error(err.Error())
			c.JSON(http.StatusBadRequest, "cant parse request")
			return
		}

		token := mocks.GetToken(usr.firstName, usr.lastName)

		c.JSON(http.StatusOK, struct{ Token string }{Token: token})
	})

	start(r)
}
