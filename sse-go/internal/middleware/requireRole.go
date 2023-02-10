package middleware

import (
	"net/http"
	"sse/internal/models"
	"sse/internal/models/roles"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequireRole(role roles.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, _ := c.Get(models.User)
		user := u.(*models.UserContext)

		if user.Role != role {
			logrus.Errorf("Required role: %s", role)
			c.AbortWithStatus(http.StatusForbidden)
		}

		c.Next()
	}
}
