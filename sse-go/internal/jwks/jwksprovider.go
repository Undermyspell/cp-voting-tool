package jwks

import (
	"context"
	"log"
	"time"

	"github.com/MicahParks/keyfunc"
)

var Jwks *keyfunc.JWKS

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

	Jwks = jwks

	go func() {
		<-teardown
		log.Println("JWKS teardown")
		cancel()
	}()
}
