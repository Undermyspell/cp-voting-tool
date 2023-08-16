package mocks

import (
	"errors"
	"sse/internal/jwks"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type JwksMock struct {
	mock.Mock
}

func NewJwks() jwks.KeyfuncProvider {
	jwksMock := new(JwksMock)
	jwksMock.On("GetKeyFunc").Return(func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("error getting keyfunc")
		}
		return []byte("my_test_secret"), nil
	})
	return jwksMock
}

func (m *JwksMock) GetKeyFunc() func(token *jwt.Token) (interface{}, error) {
	args := m.Called()
	return args.Get(0).(func(token *jwt.Token) (interface{}, error))
}
