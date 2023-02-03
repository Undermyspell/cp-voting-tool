package mocks

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/mock"
)

type JwksMock struct {
	mock.Mock
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
	claims["user"] = "Hoodini Magician"
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		log.Println("Signing error")
	}

	return tokenString
}
