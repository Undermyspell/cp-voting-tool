package middleware

import (
	"net/http"
	"voting/internal/models/roles"
	"voting/shared/shared_models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequireRole(requiredRoles ...roles.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, _ := c.Get(shared_models.User)
		user := u.(*shared_models.UserContext)

		hasRequiredRole := false
		for _, r := range requiredRoles {
			if r == user.Role {
				hasRequiredRole = true
				break
			}
		}

		if !hasRequiredRole {
			logrus.Errorf("Required role: %s", requiredRoles)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
