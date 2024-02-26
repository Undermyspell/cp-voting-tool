package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"voting/internal/broker"
	"voting/internal/env"
	"voting/internal/jwks"
	"voting/internal/middleware"
	"voting/internal/models"
	"voting/internal/models/roles"
	"voting/internal/votingstorage"
	questionService "voting/services/question"
	userService "voting/services/user"

	_ "voting/docs"

	"github.com/centrifugal/centrifuge"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var initQuestionService = func(broker broker.Broker, votingStorage votingstorage.VotingStorage) questionService.QuestionService {
	return questionService.NewBrokered(broker, votingStorage)
}

var start = func(r *gin.Engine) {
	r.Run(":3333")
}
var r *gin.Engine

type contextKey int

var ginContextKey contextKey

// GinContextToContextMiddleware - at the resolver level we only have access
// to context.Context inside centrifuge, but we need the gin context. So we
// create a gin middleware to add its context to the context.Context used by
// centrifuge websocket server.
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GinContextFromContext - we recover the gin context from the context.Context
// struct where we added it just above
func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(ginContextKey)
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}
	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func handleLog(e centrifuge.LogEntry) {
	logrus.Infof("%s: %v", e.Message, e.Fields)
}

type connectData struct {
	Email string `json:"email"`
}

// Finally we can use gin context in the auth middleware of centrifuge.
func authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// We get gin ctx from context.Context struct.
		_, err := GinContextFromContext(ctx)
		if err != nil {
			logrus.Infof("Failed to retrieve gin context")
			logrus.Info(err.Error())
			return
		}

		// And now we can access gin session.
		username := "test"
		if username != "" {
			fmt.Printf("Successful websocket auth for user %s\n", username)
		} else {
			fmt.Printf("Failed websocket auth for user %s\n", username)
			return
		}
		newCtx := centrifuge.SetCredentials(ctx, &centrifuge.Credentials{
			UserID: username,
		})
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})
}

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
	node, _ := centrifuge.New(centrifuge.Config{
		LogLevel:   centrifuge.LogLevelDebug,
		LogHandler: handleLog,
	})

	node.OnConnecting(func(ctx context.Context, event centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {

		logrus.Infof("client connecting with token: %s", event.Token)
		// return centrifuge.ConnectReply{}, centrifuge.ErrorUnauthorized

		// Let's include user email into connect reply, so we can display username in chat.
		// This is an optional step, actually.
		cred, ok := centrifuge.GetCredentials(ctx)
		if !ok {
			logrus.Error("THIS IS NOT GOOD")
			return centrifuge.ConnectReply{}, centrifuge.DisconnectServerError
		}
		data, _ := json.Marshal(connectData{
			Email: cred.UserID,
		})
		return centrifuge.ConnectReply{
			Data: data,
		}, nil
	})

	node.OnConnect(func(client *centrifuge.Client) {
		transport := client.Transport()
		logrus.Infof("🟩 user %s connected via %s.", client.UserID(), transport.Name())

		client.OnRefresh(func(e centrifuge.RefreshEvent, cb centrifuge.RefreshCallback) {
			logrus.Infof("user %s connection is going to expire, refreshing", client.UserID())
			cb(centrifuge.RefreshReply{
				ExpireAt: time.Now().Unix() + 10,
			}, nil)
		})

		client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			logrus.Infof("🟨 user %s subscribes on %s", client.UserID(), e.Channel)
			cb(centrifuge.SubscribeReply{
				Options: centrifuge.SubscribeOptions{
					EmitPresence: true,
				},
			}, nil)
		})

		client.OnUnsubscribe(func(e centrifuge.UnsubscribeEvent) {
			logrus.Infof("🟦 user %s unsubscribed from %s", client.UserID(), e.Channel)
		})

		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			logrus.Infof("🟥 user %s disconnected, disconnect: %s", client.UserID(), e.Disconnect)
		})
	})

	// We also start a separate goroutine for centrifuge itself, since we
	// still need to run gin web server.
	go func() {
		if err := node.Run(); err != nil {
			logrus.Fatal(err)
		}
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			logrus.Info("Logging to stdout")

			q := models.Question{
				Id:          "123",
				Text:        "Centrifugo!",
				Votes:       23,
				Answered:    false,
				Voted:       false,
				Anonymous:   true,
				CreatorName: "John Doe",
				CreatorHash: "ysd1123a",
			}

			data, err := json.Marshal(q)
			if err != nil {
				logrus.Error("Error marshaling:", err)
				return
			}

			node.Publish("voting", data)
		}
	}()

	env.Init()
	jwks.Init()

	r = gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	var votingStorage votingstorage.VotingStorage
	if env.Env.VotingStorageInMemory {
		votingStorage = votingstorage.NewInMemory()
	} else {
		logrus.Info("WE USE REDIS")
		votingStorage = votingstorage.NewRedis()
	}

	broker := broker.New()
	questionService := initQuestionService(broker, votingStorage)
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
	r.Use(GinContextToContextMiddleware())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/connection/websocket", gin.WrapH(authMiddleware(centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{
			CheckOrigin: func(r *http.Request) bool {
				originHeader := r.Header.Get("Origin")
				if originHeader == "" {
					return true
				}
				return originHeader == "http://localhost:5173"
			},
		}))))
		v1.GET("/events", middleware.RequireAuth(jwks.GetProvider()), broker.SseStream)
		q := v1.Group("/question", middleware.RequireAuth(jwks.GetProvider()))
		q.PUT("/answer/:id", middleware.RequireRole(roles.SessionAdmin, roles.Admin), questionService.Answer)
		q.POST("/new", questionService.Add)
		q.PUT("/upvote/:id", questionService.Upvote)
		q.PUT("/undovote/:id", questionService.UndoVote)
		q.PUT("/update", questionService.Update)
		q.DELETE("/delete/:id", questionService.Delete)

		s := q.Group("/session", middleware.RequireAuth(jwks.GetProvider()))
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
