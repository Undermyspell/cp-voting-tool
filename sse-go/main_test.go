package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sse/internal/jwks"
	"sse/internal/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type QuestionApiTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *QuestionApiTestSuite) SetupSuite() {
	jwksMock := new(mocks.JwksMock)
	jwksMock.On("GetKeyFunc").Return(func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("error getting keyfunc")
		}
		return []byte("my_test_secret"), nil
	})

	initJwks = func() jwks.KeyfuncProvider { return jwksMock }
	start = func(r *gin.Engine) {}

	main()

	suite.router = r
}

func (suite *QuestionApiTestSuite) TestNewQuestion_UNAUTHORIZED_401() {
	w := httptest.NewRecorder()

	jsonData := []byte(`{
		"text": "test question?"
	}`)

	req, _ := http.NewRequest("POST", "/question/new", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), 401, w.Code)
}

func (suite *QuestionApiTestSuite) TestNewQuestion_OK_200() {
	w := httptest.NewRecorder()

	token := mocks.GetToken()

	jsonData := []byte(`{
		"text": "test question?"
	}`)

	req, _ := http.NewRequest("POST", "/question/new", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), 200, w.Code)
}

func TestQuestionApiSuite(t *testing.T) {
	suite.Run(t, new(QuestionApiTestSuite))
}
