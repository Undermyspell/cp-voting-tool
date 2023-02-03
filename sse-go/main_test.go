package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sse/internal/jwks"
	"sse/internal/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type QuestionApiTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *QuestionApiTestSuite) SetupSuite() {
	jwksMock := mocks.NewJwks()
	initJwks = func() jwks.KeyfuncProvider { return jwksMock }
	start = func(r *gin.Engine) {}

	main()

	suite.router = r
}

func (suite *QuestionApiTestSuite) TestApi_UNAUTHORIZED_401() {
	type test struct {
		name       string
		httpMethod string
		path       string
	}

	tests := []test{
		{"Question_New_UNAUTHORIZED_401", http.MethodPost, "/question/new"},
		{"Question_Upvote_UNAUTHORIZED_401", http.MethodPut, "/question/upvote"},
		{"Question_Answer_UNAUTHORIZED_401", http.MethodPut, "/question/answer"},
		{"Question_Rest_UNAUTHORIZED_401", http.MethodPost, "/question/reset"},
		{"Events_UNAUTHORIZED_401", http.MethodGet, "/events"},
	}

	for _, test := range tests {
		suite.T().Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.httpMethod, test.path, nil)
			req.Header.Set("Content-Type", "application/json")
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
		})
	}
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

func (suite *QuestionApiTestSuite) TestUpvoteQuestion_NOTFOUND_404() {
	w := httptest.NewRecorder()

	token := mocks.GetToken()

	jsonData := []byte(`{
		"id": "invalid"
	}`)

	req, _ := http.NewRequest("PUT", "/question/upvote", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), 404, w.Code)
}

func (suite *QuestionApiTestSuite) TestAnswerQuestion_NOTFOUND_404() {
	w := httptest.NewRecorder()

	token := mocks.GetToken()

	jsonData := []byte(`{
		"id": "invalid"
	}`)

	req, _ := http.NewRequest("PUT", "/question/answer", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), 404, w.Code)
}

func TestQuestionApiSuite(t *testing.T) {
	suite.Run(t, new(QuestionApiTestSuite))
}
