package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sse/dtos"
	"sse/internal/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type QuestionApiTestSuite struct {
	suite.Suite
	router    *gin.Engine
	apiPrefix string
}

func (suite *QuestionApiTestSuite) SetupSuite() {
	os.Setenv("USE_MOCK_JWKS", "true")
	start = func(r *gin.Engine) {}

	main()

	suite.router = r
	suite.apiPrefix = "/api/v1"
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
			req, _ := http.NewRequest(test.httpMethod, fmt.Sprintf("%s%s", suite.apiPrefix, test.path), nil)
			req.Header.Set("Content-Type", "application/json")
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
		})
	}
}

func (suite *QuestionApiTestSuite) TestNewQuestion_OK_200() {
	w := httptest.NewRecorder()

	token := mocks.GetToken("test", "tester")

	jsonData := []byte(`{
		"text": "test question?"
	}`)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/question/new", suite.apiPrefix), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *QuestionApiTestSuite) TestUpvoteQuestion_NOTFOUND_404() {
	w := httptest.NewRecorder()

	token := mocks.GetToken("test", "tester")

	jsonData := []byte(`{
		"id": "invalid"
	}`)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/question/upvote", suite.apiPrefix), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *QuestionApiTestSuite) TestAnswerQuestion_NOTFOUND_404() {
	w := httptest.NewRecorder()

	token := mocks.GetToken("test", "tester")

	jsonData := []byte(`{
		"id": "invalid"
	}`)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/question/answer", suite.apiPrefix), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *QuestionApiTestSuite) TestUpvoteQuestion_NOTACCEPTABLE_406() {
	w := httptest.NewRecorder()

	token := mocks.GetToken("test", "tester")

	jsonData := dtos.NewQuestionDto{Text: "new question"}
	newQuestion, _ := json.Marshal(jsonData)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/question/new", suite.apiPrefix), bytes.NewBuffer(newQuestion))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, req)

	reql, _ := http.NewRequest("GET", fmt.Sprintf("%s/question/session", suite.apiPrefix), nil)
	reql.Header.Set("Content-Type", "application/json")
	reql.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, reql)

	var questionList []dtos.QuestionDto
	body, _ := io.ReadAll(w.Body)
	json.Unmarshal(body, &questionList)

	reqv, _ := http.NewRequest("PUT", fmt.Sprintf("%s/question/upvote/%s", suite.apiPrefix, questionList[0].Id), nil)
	reqv.Header.Set("Content-Type", "application/json")
	reqv.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, reqv)

	w2 := httptest.NewRecorder()
	reqv2, _ := http.NewRequest("PUT", fmt.Sprintf("%s/question/upvote/%s", suite.apiPrefix, questionList[0].Id), nil)
	reqv2.Header.Set("Content-Type", "application/json")
	reqv2.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w2, reqv2)

	assert.Equal(suite.T(), http.StatusNotAcceptable, w2.Code)
}

func TestQuestionApiSuite(t *testing.T) {
	suite.Run(t, new(QuestionApiTestSuite))
}
