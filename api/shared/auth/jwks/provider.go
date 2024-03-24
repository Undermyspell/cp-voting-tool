package jwks

import (
	"github.com/golang-jwt/jwt/v5"
)

type KeyfuncProvider interface {
	GetKeyFunc() func(token *jwt.Token) (interface{}, error)
}
