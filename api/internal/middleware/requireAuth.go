package middleware

import (
	"net/http"
	"voting/internal/models"

	"github.com/centrifugal/centrifuge"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GinRequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext, statusCode := models.GetUserContextFromGinContext(c)

		if statusCode > 0 {
			c.AbortWithStatus(statusCode)
			return
		}

		c.Set(models.User, userContext)

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
