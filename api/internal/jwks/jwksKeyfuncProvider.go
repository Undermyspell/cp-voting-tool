package jwks

import (
	"context"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

type JwksKeyfuncProvider struct {
	jwksKeyFunc func(token *jwt.Token) (interface{}, error)
}

func (provider *JwksKeyfuncProvider) GetKeyFunc() func(token *jwt.Token) (interface{}, error) {
	return provider.jwksKeyFunc
}

func Init() *JwksKeyfuncProvider {
	jwksURL := "https://login.microsoftonline.com/ff930651-2670-491e-9a70-7847e7fbf8b7/discovery/v2.0/keys"

	options := keyfunc.Options{
		Ctx: context.TODO(),
		RefreshErrorHandler: func(err error) {
			logrus.Errorf("There was an error with the jwt.Keyfunc\nError: %s", err.Error())
		},
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	}

	jwks, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		logrus.Fatalf("Failed to create JWKS from resource at the given URL.\nError: %s", err.Error())
	}

	jwksProvider := &JwksKeyfuncProvider{
		jwksKeyFunc: jwks.Keyfunc,
	}

	return jwksProvider
}
