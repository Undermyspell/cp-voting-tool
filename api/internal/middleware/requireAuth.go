package middleware

import (
	"errors"
	"net/http"
	"sse/internal/jwks"
	"sse/internal/models"
	"sse/internal/models/roles"
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
			logrus.Errorf("Failed to parse the JWT.\nError: %s", err.Error())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			logrus.Error("The token signature could not be verified.")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userContext, err := getUserContext(token)

		if err != nil {
			logrus.Error("Failed to parse required claims from token")
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		c.Set(models.User, userContext)

		c.Next()
	}
}

func getUserContext(token *jwt.Token) (*models.UserContext, error) {
	name, okName := token.Claims.(jwt.MapClaims)["name"]
	email, okEmail := token.Claims.(jwt.MapClaims)["preferred_username"]

	if !okEmail || !okName {
		return new(models.UserContext), errors.New("claims are not valid")
	}

	userRoles, okRole := token.Claims.(jwt.MapClaims)["roles"]
	role := roles.Contributor

	if okRole {
		t := userRoles.([]interface{})
		role = roles.Role(t[0].(string))
	}

	return &models.UserContext{Name: name.(string), Email: email.(string), Role: role}, nil
}
