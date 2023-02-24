package mocks

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

func GetUserToken(firstName, lastName string) string {
	sampleSecretKey := []byte("my_test_secret")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().UTC().Add(time.Second * 3600).Unix()
	claims["name"] = fmt.Sprintf("%s %s", firstName, lastName)
	claims["preferred_username"] = fmt.Sprintf("%s %s", firstName, lastName)
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		logrus.Fatal("Signing error")
	}

	return tokenString
}

func GetAdminUserToken(firstName, lastName string) string {
	sampleSecretKey := []byte("my_test_secret")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().UTC().Add(time.Second * 3600).Unix()
	claims["name"] = fmt.Sprintf("%s %s", firstName, lastName)
	claims["preferred_username"] = fmt.Sprintf("%s.%s@mock.com", firstName, lastName)
	claims["roles"] = []string{"admin"}
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		logrus.Fatal("Signing error")
	}

	return tokenString
}

func GetSessionAdminUserToken(firstName, lastName string) string {
	sampleSecretKey := []byte("my_test_secret")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().UTC().Add(time.Second * 3600).Unix()
	claims["name"] = fmt.Sprintf("%s %s", firstName, lastName)
	claims["preferred_username"] = fmt.Sprintf("%s.%s@mock.com", firstName, lastName)
	claims["roles"] = []string{"session_admin"}
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		logrus.Fatal("Signing error")
	}

	return tokenString
}
