package mocks

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type JwksMock struct {
	mock.Mock
}

func NewJwks() *JwksMock {
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

func GetToken() string {
	sampleSecretKey := []byte("my_test_secret")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().UTC().Add(time.Second * 3600).Unix()
	claims["name"] = "Test Tester"
	claims["email"] = "testuser@test.com"
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		logrus.Fatal("Signing error")
	}

	return tokenString
}
