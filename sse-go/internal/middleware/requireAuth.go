package middleware

import (
	"net/http"
	"sse/internal/jwks"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

func RequireAuth(keyfuncProvider jwks.KeyfuncProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")

		splitted := strings.Split(bearerToken, " ")

		if len(splitted) != 2 && splitted[0] != "Bearer" {
			logrus.Error("Not a valid bearer token")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		jwtB64 := splitted[1]
		token, err := jwt.Parse(jwtB64, keyfuncProvider.GetKeyFunc())

		if err != nil {
			logrus.Error("Failed to parse the JWT.\nError: %s", err.Error())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			logrus.Error("The token is not valid.")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
