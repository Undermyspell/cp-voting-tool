package middleware

import (
	"net/http"
	"time"
	shared_models "voting/shared/models"

	"github.com/centrifugal/centrifuge"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GinRequireJwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext, statusCode := shared_models.GetUserContextFromGinContext(c)

		if statusCode > 0 {
			c.AbortWithStatus(statusCode)
			return
		}

		c.Set(shared_models.User, userContext)

		c.Next()
	}
}

func GinRequireCookieAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/oauth2/callback" {
			c.Next()
			return
		}

		session := sessions.Default(c)

		auth, ok := session.Get("authenticated").(bool)
		token, _ := session.Get("token").(string)
		tokenExpiry, _ := session.Get("tokenExpiry").(int64)

		logrus.Info(time.Now().Unix())
		logrus.Info(tokenExpiry)

		logrus.Printf("Session 'authenticated': %v, ok: %v", auth, ok)
		if !ok || !auth || token == "" || time.Now().Unix() > tokenExpiry {
			logrus.Warn("User is not authenticated, redirecting to login")
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}

func CentrifugeAnonymousAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// We get gin ctx from context.Context struct.
		_, err := GinContextFromContext(ctx)
		if err != nil {
			logrus.Infof("Failed to retrieve gin context")
			logrus.Info(err.Error())
			return
		}

		newCtx := centrifuge.SetCredentials(ctx, &centrifuge.Credentials{
			UserID: "",
		})
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})
}
