package jwks

import (
	"context"
	"log"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

type KeyfuncProvider interface {
	GetKeyFunc() func(token *jwt.Token) (interface{}, error)
}

type JwksKeyfuncProvider struct {
	jwksKeyFunc func(token *jwt.Token) (interface{}, error)
}

func (provider *JwksKeyfuncProvider) GetKeyFunc() func(token *jwt.Token) (interface{}, error) {
	return provider.jwksKeyFunc
}

func Init() *JwksKeyfuncProvider {
	// ctx, cancel := context.WithCancel(context.Background())
	jwksURL := ""

	options := keyfunc.Options{
		Ctx: context.TODO(),
		RefreshErrorHandler: func(err error) {
			log.Printf("There was an error with the jwt.Keyfunc\nError: %s", err.Error())
		},
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	}

	jwks, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		log.Fatalf("Failed to create JWKS from resource at the given URL.\nError: %s", err.Error())
	}

	jwksProvider := &JwksKeyfuncProvider{
		jwksKeyFunc: jwks.Keyfunc,
	}

	// go func() {
	// 	<-teardown
	// 	log.Println("JWKS teardown")
	// 	cancel()
	// }()

	return jwksProvider
}
