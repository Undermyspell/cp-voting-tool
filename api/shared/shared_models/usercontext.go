package shared_models

import (
	"crypto/hmac"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"voting/internal/jwks"
	"voting/internal/models/roles"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
)

type UserContext struct {
	Name  string
	Email string
	Role  roles.Role
}

func (userContext *UserContext) GetHash(secret string) string {
	h := hmac.New(sha3.New256, []byte(secret))
	marshalled, _ := json.Marshal(userContext)
	h.Write(marshalled)
	hash := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprint(hash)
}

const User string = "user"

func GetUserContextFromToken(jwtToken string) (*UserContext, error) {
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

	return &UserContext{Name: name.(string), Email: email.(string), Role: role}, nil
}

func GetUserContextFromGinContext(c *gin.Context) (*UserContext, int) {
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
