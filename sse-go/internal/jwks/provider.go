package jwks

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

type KeyfuncProvider interface {
	GetKeyFunc() func(token *jwt.Token) (interface{}, error)
}

type JwksKeyfuncProvider struct{}

func (provider *JwksKeyfuncProvider) GetKeyFunc() func(token *jwt.Token) (interface{}, error) {
	return jwksKeyFunc
}

type MockKeyfuncProvider struct{}

func (provider *MockKeyfuncProvider) GetKeyFunc() func(token *jwt.Token) (interface{}, error) {
	return mockKeyFunc
}

var jwksKeyFunc func(token *jwt.Token) (interface{}, error)
var mockKeyFunc func(token *jwt.Token) (interface{}, error) = func(token *jwt.Token) (interface{}, error) {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, errors.New("error getting keyfunc")
	}
	return []byte("my_test_secret"), nil
}

func Init(teardown chan string) {
	ctx, cancel := context.WithCancel(context.Background())
	jwksURL := ""

	options := keyfunc.Options{
		Ctx: ctx,
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

	jwksKeyFunc = jwks.Keyfunc

	go func() {
		<-teardown
		log.Println("JWKS teardown")
		cancel()
	}()
}
