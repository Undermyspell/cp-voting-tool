package middleware

import (
	"errors"
	"net/http"
	"strings"
	"voting/internal/jwks"
	"voting/internal/models"
	"voting/internal/models/roles"

	"github.com/centrifugal/centrifuge"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func GinRequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext, statusCode := getUserContext(c)

		if statusCode > 0 {
			c.AbortWithStatus(statusCode)
			return
		}

		c.Set(models.User, userContext)

		c.Next()
	}
}

// Finally we can use gin context in the auth middleware of centrifuge.
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

func GetUserContextFromToken(jwtToken string) (*models.UserContext, error) {
	token, err := jwt.Parse(jwtToken, jwks.GetProvider().GetKeyFunc())

	if err != nil {
		logrus.Errorf("Failed to parse the JWT.\nError: %s", err.Error())
		return nil, err
	}

	if !token.Valid {
		logrus.Error("The token signature could not be verified.")
		return nil, err
	}

	name, okName := token.Claims.(jwt.MapClaims)["name"]
	email, okEmail := token.Claims.(jwt.MapClaims)["preferred_username"]

	if !okEmail || !okName {
		logrus.Error("Failed to parse required claims from token")
		return nil, errors.New("claims are not valid")
	}

	userRoles, okRole := token.Claims.(jwt.MapClaims)["roles"]
	role := roles.Contributor

	if okRole {
		t := userRoles.([]interface{})
		role = roles.Role(t[0].(string))
	}

	return &models.UserContext{Name: name.(string), Email: email.(string), Role: role}, nil
}

func getUserContext(c *gin.Context) (*models.UserContext, int) {
	bearerToken := c.GetHeader("Authorization")

	splitted := strings.Split(bearerToken, " ")

	if len(splitted) != 2 && splitted[0] != "Bearer" {
		logrus.Error("Not a valid bearer token")
		return nil, http.StatusUnauthorized
	}

	jwtB64 := splitted[1]
	userContext, err := GetUserContextFromToken(jwtB64)

	if err != nil {
		return nil, http.StatusUnauthorized
	}

	return userContext, 0
}
