package jwks

import (
	"voting/internal/env"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

var jwksProvider KeyfuncProvider

type JwksKeyfuncProvider struct {
	jwksKeyFunc func(token *jwt.Token) (interface{}, error)
}

func (provider *JwksKeyfuncProvider) GetKeyFunc() func(token *jwt.Token) (interface{}, error) {
	return provider.jwksKeyFunc
}

func GetProvider() KeyfuncProvider {
	return jwksProvider
}

func create() {
	jwksURL := env.Env.JwksUrl

	jwks, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		logrus.Fatalf("Failed to create JWKS from resource at the given URL.\nError: %s", err.Error())
	}

	jwksProvider = &JwksKeyfuncProvider{
		jwksKeyFunc: jwks.Keyfunc,
	}
}
