package jwks

import (
	"github.com/golang-jwt/jwt/v4"
)

type KeyfuncProvider interface {
	GetKeyFunc() func(token *jwt.Token) (interface{}, error)
}
