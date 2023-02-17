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
	"golang.org/x/exp/slices"
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
		{"Question_Upvote_UNAUTHORIZED_401", http.MethodPut, "/question/upvote/question1"},
		{"Question_Answer_UNAUTHORIZED_401", http.MethodPut, "/question/answer/question1"},
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

func (suite *QuestionApiTestSuite) TestGetSession_OK_200_CREATOR_SHOWN_ONLY_FOR_OWNED_AND_NO_ANONYMOUS_QUESTIONS() {
	w := httptest.NewRecorder()

	tokenUser_Foo := mocks.GetToken("Foo", "Foo_Tester")
	tokenUser_Bar := mocks.GetToken("Bar", "Bar_Tester")

	newQuestion_FOO_Q1 := dtos.NewQuestionDto{Text: "Foo Question1", Anonymous: false}
	newQuestion_FOO_Q2 := dtos.NewQuestionDto{Text: "Foo Question2 anonynmous", Anonymous: true}
	newQuestion_BAR_Q1 := dtos.NewQuestionDto{Text: "Bar Question1", Anonymous: false}
	newQuestion_BAR_Q2 := dtos.NewQuestionDto{Text: "Bar Question2 anonynmous", Anonymous: true}

	postNewQuestion(suite, w, newQuestion_FOO_Q1, tokenUser_Foo)
	postNewQuestion(suite, w, newQuestion_FOO_Q2, tokenUser_Foo)
	postNewQuestion(suite, w, newQuestion_BAR_Q1, tokenUser_Bar)
	postNewQuestion(suite, w, newQuestion_BAR_Q2, tokenUser_Bar)

	reql, _ := http.NewRequest("GET", fmt.Sprintf("%s/question/session", suite.apiPrefix), nil)
	reql.Header.Set("Content-Type", "application/json")
	reql.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenUser_Foo))
	suite.router.ServeHTTP(w, reql)

	var questionList []dtos.QuestionDto
	body, _ := io.ReadAll(w.Body)
	json.Unmarshal(body, &questionList)

	question_FOO_Q1 := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Foo Question1" })]
	question_FOO_Q2 := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Foo Question2 anonynmous" })]
	question_BAR_Q1 := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Bar Question1" })]
	question_BAR_Q2 := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Bar Question2 anonynmous" })]

	assert.Equal(suite.T(), true, question_FOO_Q1.Owned)
	assert.Equal(suite.T(), false, question_FOO_Q1.Anonymous)
	assert.Equal(suite.T(), "Foo Foo_Tester", question_FOO_Q1.Creator)

	assert.Equal(suite.T(), true, question_FOO_Q2.Owned)
	assert.Equal(suite.T(), true, question_FOO_Q2.Anonymous)
	assert.Equal(suite.T(), "Foo Foo_Tester", question_FOO_Q2.Creator)

	assert.Equal(suite.T(), false, question_BAR_Q1.Owned)
	assert.Equal(suite.T(), false, question_BAR_Q1.Anonymous)
	assert.Equal(suite.T(), "Bar Bar_Tester", question_BAR_Q1.Creator)

	assert.Equal(suite.T(), false, question_BAR_Q2.Owned)
	assert.Equal(suite.T(), true, question_BAR_Q2.Anonymous)
	assert.Equal(suite.T(), "", question_BAR_Q2.Creator)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func TestQuestionApiSuite(t *testing.T) {
	suite.Run(t, new(QuestionApiTestSuite))
}

func postNewQuestion(suite *QuestionApiTestSuite, w *httptest.ResponseRecorder, jsonDataQ1 dtos.NewQuestionDto, token string) {
	newQuestion, _ := json.Marshal(jsonDataQ1)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/question/new", suite.apiPrefix), bytes.NewBuffer(newQuestion))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, req)
}
