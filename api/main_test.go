package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"sync"
	"testing"
	"time"
	"voting/dtos"
	"voting/internal/events"
	"voting/internal/mocks"

	"github.com/centrifugal/centrifuge-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CentrifugeTestClient struct {
	client           *centrifuge.Client
	receivedMessages []events.Event
}

type QuestionApiTestSuite struct {
	suite.Suite
	router                 *gin.Engine
	apiPrefix              string
	tokenUser_Foo          string
	tokenUser_Bar          string
	tokenUser_Admin        string
	tokenUser_SessionAdmin string
	centrifugeClientFoo    CentrifugeTestClient
	centrifugeClientBar    CentrifugeTestClient
}

func (suite *QuestionApiTestSuite) SetupSuite() {
	os.Setenv("USE_MOCK_JWKS", "true")
	os.Setenv("GENERATE_REDIS_STORAGE_ROOT_KEY", "true")
	os.Setenv("JWKS_URL", "https://test/discovery/v2.0/keys")
	os.Setenv("VOTING_STORAGE_IN_MEMORY", "true")
	os.Setenv("VOTING_STORAGE_IN_MEMORY", "true")
	os.Setenv("ALLOWED_ORIGIN", "http://localhost:5173")

	start = func(r *gin.Engine) {
		go func() {
			r.Run(":3333")

		}()
	}

	main()

	suite.router = r
	suite.apiPrefix = "/api/v1"
	suite.tokenUser_Foo = mocks.GetUserToken("Foo", "Foo_Tester")
	suite.tokenUser_Bar = mocks.GetUserToken("Bar", "Bar_Tester")
	suite.tokenUser_Admin = mocks.GetAdminUserToken("Admin", "Admin_Tester")
	suite.tokenUser_SessionAdmin = mocks.GetSessionAdminUserToken("SessionAdmin", "Session_Admin_Tester")

	initCentrifuge(suite)

	suite.centrifugeClientFoo.client.OnMessage(func(me centrifuge.MessageEvent) {

		var event events.Event

		json.Unmarshal([]byte(me.Data), &event)
		suite.centrifugeClientFoo.receivedMessages = append(suite.centrifugeClientFoo.receivedMessages, event)
	})

	suite.centrifugeClientBar.client.OnMessage(func(me centrifuge.MessageEvent) {
		var event events.Event
		json.Unmarshal([]byte(me.Data), &event)

		suite.centrifugeClientBar.receivedMessages = append(suite.centrifugeClientBar.receivedMessages, event)
	})
}

func (suite *QuestionApiTestSuite) TearDownTest() {
	stopSession(suite)
}

func (suite *QuestionApiTestSuite) SetupTest() {
	clearCentrifugeClientsMessages(suite)
	startSession(suite)
}

func (suite *QuestionApiTestSuite) TearDownSuite() {
	suite.centrifugeClientFoo.client.Close()
	suite.centrifugeClientBar.client.Close()
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
		{"Question_Undovote_UNAUTHORIZED_401", http.MethodPut, "/question/undovote/question1"},
		{"Question_Answer_UNAUTHORIZED_401", http.MethodPut, "/question/answer/question1"},
		{"Question_Rest_UNAUTHORIZED_401", http.MethodGet, "/question/session"},
		{"Question_Rest_UNAUTHORIZED_401", http.MethodPost, "/question/session/start"},
		{"Question_Rest_UNAUTHORIZED_401", http.MethodPost, "/question/session/stop"},
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

func (suite *QuestionApiTestSuite) TestApi_NOTACCEPTABLE_406_WHEN_NO_SESSION_RUNNING() {
	type test struct {
		name       string
		httpMethod string
		path       string
		payload    *bytes.Buffer
	}

	stopSession(suite)

	token := suite.tokenUser_Admin

	tests := []test{
		{"Question_New_NOTACCEPTABLE_406", http.MethodPost, "/question/new", bytes.NewBuffer([]byte(`{"text": "test question?" }`))},
		{"Question_Upvote_NOTACCEPTABLE_406", http.MethodPut, "/question/upvote/question1", nil},
		{"Question_Undovote_NOTACCEPTABLE_406", http.MethodPut, "/question/undovote/question1", nil},
		{"Question_Answer_NOTACCEPTABLE_406", http.MethodPut, "/question/answer/question1", nil},
		{"Question_Update_NOTACCEPTABLE_406", http.MethodPut, "/question/update", bytes.NewBuffer([]byte(`{"id": "1","text": "test question?" }`))},
		{"Question_Delete_NOTACCEPTABLE_406", http.MethodDelete, "/question/delete/question1", nil},
		{"Question_GetSession_NOTACCEPTABLE_406", http.MethodGet, "/question/session", nil},
	}

	for _, test := range tests {
		suite.T().Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			var req *http.Request
			if test.payload == nil {
				req, _ = http.NewRequest(test.httpMethod, fmt.Sprintf("%s%s", suite.apiPrefix, test.path), nil)
			} else {
				req, _ = http.NewRequest(test.httpMethod, fmt.Sprintf("%s%s", suite.apiPrefix, test.path), test.payload)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusNotAcceptable, w.Code)
		})
	}
}

func (suite *QuestionApiTestSuite) TestNewQuestion_OK_200() {
	w := httptest.NewRecorder()

	token := suite.tokenUser_Foo
	newQuestion := dtos.NewQuestionDto{Text: "Foo Question", Anonymous: false}

	postNewQuestion(suite, w, newQuestion, token)

	questionList := getSession(suite, w, token)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), 1, questionList[0].Votes)
	assert.Equal(suite.T(), true, questionList[0].Voted)
}

func (suite *QuestionApiTestSuite) TestNewQuestion_EVENTS_RECEIVED_NEW_QUESTION() {
	w := httptest.NewRecorder()

	token := suite.tokenUser_Foo
	newQuestion := dtos.NewQuestionDto{Text: "This is a new question", Anonymous: false}

	postNewQuestion(suite, w, newQuestion, token)

	assert.EventuallyWithT(suite.T(), func(c *assert.CollectT) {
		questionCreatedEventFoo, err := findEvent[events.QuestionCreated](suite.centrifugeClientFoo, events.NEW_QUESTION)
		assert.Nil(c, err)
		assert.Equal(c, "This is a new question", questionCreatedEventFoo.Text)
		assert.Equal(c, false, questionCreatedEventFoo.Answered)
		assert.Equal(c, false, questionCreatedEventFoo.Anonymous)
		assert.Equal(c, "Foo Foo_Tester", questionCreatedEventFoo.Creator)
		assert.Equal(c, true, questionCreatedEventFoo.Owned)
		assert.Equal(c, true, questionCreatedEventFoo.Voted)
		assert.Equal(c, 1, questionCreatedEventFoo.Votes)

		questionCreatedEventBar, err1 := findEvent[events.QuestionCreated](suite.centrifugeClientBar, events.NEW_QUESTION)
		assert.Nil(c, err1)
		assert.Equal(c, "This is a new question", questionCreatedEventBar.Text)
		assert.Equal(c, false, questionCreatedEventBar.Answered)
		assert.Equal(c, false, questionCreatedEventBar.Anonymous)
		assert.Equal(c, "Foo Foo_Tester", questionCreatedEventBar.Creator)
		assert.Equal(c, false, questionCreatedEventBar.Owned)
		assert.Equal(c, false, questionCreatedEventBar.Voted)
		assert.Equal(c, 1, questionCreatedEventBar.Votes)
	}, time.Second*3, time.Second*1, "New question event has not been received")
}

func (suite *QuestionApiTestSuite) TestNewQuestion_404_WHEN_TEXT_LENTGH_OVER_MAX_LENGTH() {
	w := httptest.NewRecorder()

	token := suite.tokenUser_Foo
	newQuestion := dtos.NewQuestionDto{Text: "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.", Anonymous: false}

	postNewQuestion(suite, w, newQuestion, token)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *QuestionApiTestSuite) TestUpvoteQuestion_NOTFOUND_404() {
	w := httptest.NewRecorder()

	token := suite.tokenUser_Foo

	upvoteQuestion(suite, w, "invalid_id", token)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *QuestionApiTestSuite) TestAnswerQuestion_NOTFOUND_404() {
	w := httptest.NewRecorder()

	token := suite.tokenUser_SessionAdmin

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/question/answer/%s", suite.apiPrefix, "invalid"), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *QuestionApiTestSuite) TestAnswerQuestion_OK_200() {
	type test struct {
		name  string
		token string
	}

	tests := []test{
		{"SESSION_ADMIN_200", suite.tokenUser_SessionAdmin},
		{"ADMIN_200", suite.tokenUser_Admin},
	}

	for _, test := range tests {
		suite.T().Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			newQuestion := dtos.NewQuestionDto{Text: "Foo Question", Anonymous: false}
			postNewQuestion(suite, w, newQuestion, suite.tokenUser_Bar)
			questionList := getSession(suite, w, suite.tokenUser_Bar)
			question_FOO_Q := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Foo Question" })]

			assert.Equal(suite.T(), "Foo Question", question_FOO_Q.Text)
			assert.Equal(suite.T(), false, question_FOO_Q.Anonymous)

			w2 := httptest.NewRecorder()
			answerQuestion(suite, w2, question_FOO_Q.Id, test.token)

			assert.Equal(suite.T(), http.StatusOK, w.Code)
		})
	}
}

func (suite *QuestionApiTestSuite) TestAnswerQuestion_EVENTS_RECEIVED_ANSWER_QUESTION() {
	w := httptest.NewRecorder()

	newQuestion := dtos.NewQuestionDto{Text: "Question to be answered", Anonymous: false}
	postNewQuestion(suite, w, newQuestion, suite.tokenUser_Foo)
	questionList := getSession(suite, w, suite.tokenUser_Foo)
	question_FOO_Q := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Question to be answered" })]

	answerQuestion(suite, w, question_FOO_Q.Id, suite.tokenUser_Admin)

	assert.EventuallyWithT(suite.T(), func(c *assert.CollectT) {
		questionAnsweredEventFoo, err := findEvent[events.QuestionAnswered](suite.centrifugeClientFoo, events.ANSWER_QUESTION)
		assert.Nil(c, err)
		assert.Equal(c, question_FOO_Q.Id, questionAnsweredEventFoo.Id)

		questionAnsweredEventBar, err1 := findEvent[events.QuestionAnswered](suite.centrifugeClientBar, events.ANSWER_QUESTION)
		assert.Nil(c, err1)
		assert.Equal(c, question_FOO_Q.Id, questionAnsweredEventBar.Id)
	}, time.Second*3, time.Second*1, "No question answered event received")
}

func (suite *QuestionApiTestSuite) TestAnswerQuestion_FORBIDDEN_403_WHEN_NO_SESSION_ADMIN_OR_ADMIN() {
	tokenUser_Foo := suite.tokenUser_Foo

	w := httptest.NewRecorder()

	answerQuestion(suite, w, "id_123", tokenUser_Foo)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

func (suite *QuestionApiTestSuite) TestUpvoteQuestion_NOTACCEPTABLE_406_WHEN_DOUBLE_VOTE_FROM_USER() {
	w := httptest.NewRecorder()

	token := suite.tokenUser_Foo
	newQuestion := dtos.NewQuestionDto{Text: "new question"}
	postNewQuestion(suite, w, newQuestion, token)

	questionList := getSession(suite, w, token)

	upvoteQuestion(suite, w, questionList[0].Id, token)

	w2 := httptest.NewRecorder()
	upvoteQuestion(suite, w2, questionList[0].Id, token)

	assert.Equal(suite.T(), http.StatusNotAcceptable, w2.Code)
}

func (suite *QuestionApiTestSuite) TestUpvoteQuestion_NOTACCEPTABLE_406_WHEN_VOTING_ANSWERED_QUESTION() {
	w := httptest.NewRecorder()

	token := suite.tokenUser_SessionAdmin

	jsonData := dtos.NewQuestionDto{Text: "new question"}

	postNewQuestion(suite, w, jsonData, token)

	questionList := getSession(suite, w, token)

	answerQuestion(suite, w, questionList[0].Id, token)

	w2 := httptest.NewRecorder()
	upvoteQuestion(suite, w2, questionList[0].Id, token)

	assert.Equal(suite.T(), http.StatusNotAcceptable, w2.Code)
}

func (suite *QuestionApiTestSuite) TestDeleteQuestion_NOTACCEPTABLE_406_WHEN_DELETING_ANSWERED_QUESTION() {
	w := httptest.NewRecorder()

	token := suite.tokenUser_SessionAdmin

	jsonData := dtos.NewQuestionDto{Text: "new question"}

	postNewQuestion(suite, w, jsonData, token)

	questionList := getSession(suite, w, token)

	answerQuestion(suite, w, questionList[0].Id, token)

	w2 := httptest.NewRecorder()
	deleteQuestion(suite, w2, questionList[0].Id, token)

	assert.Equal(suite.T(), http.StatusNotAcceptable, w2.Code)
}

func (suite *QuestionApiTestSuite) TestUpdateQuestion_NOTACCEPTABLE_406_WHEN_UPDATING_ANSWERED_QUESTION() {
	w := httptest.NewRecorder()

	token := suite.tokenUser_SessionAdmin

	jsonData := dtos.NewQuestionDto{Text: "new question"}

	postNewQuestion(suite, w, jsonData, token)

	questionList := getSession(suite, w, token)

	answerQuestion(suite, w, questionList[0].Id, token)

	w2 := httptest.NewRecorder()
	updateQuestionDto := dtos.UpdateQuestionDto{Id: questionList[0].Id, Text: "Updated Question", Anonymous: true}
	putUpdateQuestion(suite, w2, updateQuestionDto, token)

	assert.Equal(suite.T(), http.StatusNotAcceptable, w2.Code)
}

func (suite *QuestionApiTestSuite) TestUpvoteQuestion_SAME_QUESTION_PARALLEL_100_SHOULD_RETURN_100() {
	w := httptest.NewRecorder()
	jsonData := dtos.NewQuestionDto{Text: "new question"}
	postNewQuestion(suite, w, jsonData, suite.tokenUser_Bar)
	questionList := getSession(suite, w, suite.tokenUser_Bar)
	questionId := questionList[0].Id

	suite.T().Run("Parallel_Question_Upvote", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 1; i <= 99; i++ {
			wg.Add(1)
			tokenUser := mocks.GetUserToken(fmt.Sprintf("User_%d", i), fmt.Sprintf("User_%d", i))
			go func(tokenUser string, w *httptest.ResponseRecorder, questionId string) {
				defer wg.Done()
				upvoteQuestion(suite, w, questionId, tokenUser)
			}(tokenUser, w, questionId)
		}
		wg.Wait()

		w1 := httptest.NewRecorder()

		questionList := getSession(suite, w1, suite.tokenUser_Bar)
		question := questionList[0]

		assert.EventuallyWithT(suite.T(), func(c *assert.CollectT) {
			assert.Equal(suite.T(), 100, question.Votes)
		}, time.Second*20, time.Second*1, "Question votes count does not match")
	})
}

func (suite *QuestionApiTestSuite) TestUpvoteQuestion_EVENTS_RECEIVED_UPVOTE_QUESTION() {
	w := httptest.NewRecorder()
	jsonData := dtos.NewQuestionDto{Text: "new question"}
	postNewQuestion(suite, w, jsonData, suite.tokenUser_Bar)
	questionList := getSession(suite, w, suite.tokenUser_Bar)
	questionId := questionList[0].Id

	upvoteQuestion(suite, w, questionId, suite.tokenUser_Foo)

	assert.EventuallyWithT(suite.T(), func(c *assert.CollectT) {
		questionUpvotedEventFoo, err := findEvent[events.QuestionUpvoted](suite.centrifugeClientFoo, events.UPVOTE_QUESTION)
		assert.Nil(c, err)
		assert.Equal(c, questionId, questionUpvotedEventFoo.Id)
		assert.Equal(c, true, questionUpvotedEventFoo.Voted)
		assert.Equal(c, 2, questionUpvotedEventFoo.Votes)

		questionUpvotedEventBar, err1 := findEvent[events.QuestionUpvoted](suite.centrifugeClientBar, events.UPVOTE_QUESTION)
		assert.Nil(c, err1)
		assert.Equal(c, questionId, questionUpvotedEventBar.Id)
		assert.Equal(c, false, questionUpvotedEventBar.Voted)
		assert.Equal(c, 2, questionUpvotedEventBar.Votes)
	}, time.Second*3, time.Second*1, "No question upvote event received")
}

func (suite *QuestionApiTestSuite) TestUndVoteQuestion_NOTACCEPTABLE_406_WHEN_USER_NOT_VOTED() {
	w := httptest.NewRecorder()

	token := suite.tokenUser_Foo
	newQuestion := dtos.NewQuestionDto{Text: "new question"}
	postNewQuestion(suite, w, newQuestion, token)

	questionList := getSession(suite, w, token)

	undoVoteQuestion(suite, w, questionList[0].Id, token)

	w2 := httptest.NewRecorder()
	undoVoteQuestion(suite, w2, questionList[0].Id, token)

	assert.Equal(suite.T(), http.StatusNotAcceptable, w2.Code)
}

func (suite *QuestionApiTestSuite) TestUndovoteQuestion_NOTACCEPTABLE_406_WHEN_UNDO_VOTE_ANSWERED_QUESTION() {
	w := httptest.NewRecorder()

	token := suite.tokenUser_SessionAdmin

	jsonData := dtos.NewQuestionDto{Text: "new question"}

	postNewQuestion(suite, w, jsonData, token)

	questionList := getSession(suite, w, token)

	answerQuestion(suite, w, questionList[0].Id, token)

	w2 := httptest.NewRecorder()
	undoVoteQuestion(suite, w2, questionList[0].Id, token)

	assert.Equal(suite.T(), http.StatusNotAcceptable, w2.Code)
}

func (suite *QuestionApiTestSuite) TestUndovoteeQuestion_NOTFOUND_404() {
	w := httptest.NewRecorder()

	token := suite.tokenUser_Foo

	undoVoteQuestion(suite, w, "invalid_id", token)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *QuestionApiTestSuite) TestUndovoteQuestion_SAME_QUESTION_PARALLEL_100_SHOULD_RETURN_0() {
	w := httptest.NewRecorder()
	jsonData := dtos.NewQuestionDto{Text: "new question"}
	postNewQuestion(suite, w, jsonData, suite.tokenUser_Bar)
	questionList := getSession(suite, w, suite.tokenUser_Bar)
	questionId := questionList[0].Id

	suite.T().Run("Parallel_Question_Upvote", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 1; i <= 99; i++ {
			wg.Add(1)
			tokenUser := mocks.GetUserToken(fmt.Sprintf("User_%d", i), fmt.Sprintf("User_%d", i))
			go func(tokenUser string, w *httptest.ResponseRecorder, questionId string) {
				defer wg.Done()
				upvoteQuestion(suite, w, questionId, tokenUser)
			}(tokenUser, w, questionId)
		}
		wg.Wait()
	})

	suite.T().Run("Parallel_Question_Undoupvote", func(t *testing.T) {
		var wg sync.WaitGroup
		undoVoteQuestion(suite, w, questionId, suite.tokenUser_Bar)
		for i := 1; i <= 99; i++ {
			wg.Add(1)
			tokenUser := mocks.GetUserToken(fmt.Sprintf("User_%d", i), fmt.Sprintf("User_%d", i))
			go func(tokenUser string, w *httptest.ResponseRecorder, questionId string) {
				defer wg.Done()
				undoVoteQuestion(suite, w, questionId, tokenUser)
			}(tokenUser, w, questionId)
		}
		wg.Wait()

		w1 := httptest.NewRecorder()

		questionList := getSession(suite, w1, suite.tokenUser_Bar)
		question := questionList[0]

		assert.EventuallyWithT(suite.T(), func(c *assert.CollectT) {
			assert.Equal(suite.T(), 0, question.Votes)
		}, time.Second*20, time.Second*1, "Question votes count is not 0")
	})
}

func (suite *QuestionApiTestSuite) TestUndovoteQuestion_EVENTS_RECEIVED_UNDO_UPVOTE_QUESTION() {
	w := httptest.NewRecorder()
	jsonData := dtos.NewQuestionDto{Text: "new question"}
	postNewQuestion(suite, w, jsonData, suite.tokenUser_Bar)
	questionList := getSession(suite, w, suite.tokenUser_Bar)
	questionId := questionList[0].Id

	upvoteQuestion(suite, w, questionId, suite.tokenUser_Foo)
	clearCentrifugeClientsMessages(suite)
	undoVoteQuestion(suite, w, questionId, suite.tokenUser_Bar)

	assert.EventuallyWithT(suite.T(), func(c *assert.CollectT) {
		questionUpvotedEventFoo, err := findEvent[events.QuestionUpvoted](suite.centrifugeClientFoo, events.UNDO_UPVOTE_QUESTION)
		assert.Nil(c, err)
		assert.Equal(c, questionId, questionUpvotedEventFoo.Id)
		assert.Equal(c, false, questionUpvotedEventFoo.Voted)
		assert.Equal(c, 1, questionUpvotedEventFoo.Votes)

		questionUpvotedEventBar, err1 := findEvent[events.QuestionUpvoted](suite.centrifugeClientBar, events.UNDO_UPVOTE_QUESTION)
		assert.Nil(c, err1)
		assert.Equal(c, questionId, questionUpvotedEventBar.Id)
		assert.Equal(c, false, questionUpvotedEventBar.Voted)
		assert.Equal(c, 1, questionUpvotedEventBar.Votes)
	}, time.Second*3, time.Second*1, "No question undovote event received")
}

func (suite *QuestionApiTestSuite) TestGetSession_OK_200_CREATOR_SHOWN_ONLY_FOR_OWNED_AND_NOT_ANONYMOUS_QUESTIONS() {
	w := httptest.NewRecorder()

	tokenUser_Foo := suite.tokenUser_Foo
	tokenUser_Bar := suite.tokenUser_Bar

	newQuestion_FOO_Q1 := dtos.NewQuestionDto{Text: "Foo Question1", Anonymous: false}
	newQuestion_FOO_Q2 := dtos.NewQuestionDto{Text: "Foo Question2 anonynmous", Anonymous: true}
	newQuestion_BAR_Q1 := dtos.NewQuestionDto{Text: "Bar Question1", Anonymous: false}
	newQuestion_BAR_Q2 := dtos.NewQuestionDto{Text: "Bar Question2 anonynmous", Anonymous: true}

	postNewQuestion(suite, w, newQuestion_FOO_Q1, tokenUser_Foo)
	postNewQuestion(suite, w, newQuestion_FOO_Q2, tokenUser_Foo)
	postNewQuestion(suite, w, newQuestion_BAR_Q1, tokenUser_Bar)
	postNewQuestion(suite, w, newQuestion_BAR_Q2, tokenUser_Bar)

	questionList := getSession(suite, w, tokenUser_Foo)

	question_FOO_Q1 := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Foo Question1" })]
	question_FOO_Q2 := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Foo Question2 anonynmous" })]
	question_BAR_Q1 := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Bar Question1" })]
	question_BAR_Q2 := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Bar Question2 anonynmous" })]

	assert.Equal(suite.T(), true, question_FOO_Q1.Owned)
	assert.Equal(suite.T(), false, question_FOO_Q1.Anonymous)
	assert.Equal(suite.T(), "Foo Foo_Tester", question_FOO_Q1.Creator)

	assert.Equal(suite.T(), true, question_FOO_Q2.Owned)
	assert.Equal(suite.T(), true, question_FOO_Q2.Anonymous)
	assert.Equal(suite.T(), "", question_FOO_Q2.Creator)

	assert.Equal(suite.T(), false, question_BAR_Q1.Owned)
	assert.Equal(suite.T(), false, question_BAR_Q1.Anonymous)
	assert.Equal(suite.T(), "Bar Bar_Tester", question_BAR_Q1.Creator)

	assert.Equal(suite.T(), false, question_BAR_Q2.Owned)
	assert.Equal(suite.T(), true, question_BAR_Q2.Anonymous)
	assert.Equal(suite.T(), "", question_BAR_Q2.Creator)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *QuestionApiTestSuite) TestUpdateQuestion_OK_200_WHEN_QUESTION_IS_OWNED() {
	w := httptest.NewRecorder()

	tokenUser_Foo := suite.tokenUser_Foo
	newQuestion := dtos.NewQuestionDto{Text: "Foo Question", Anonymous: false}
	postNewQuestion(suite, w, newQuestion, tokenUser_Foo)
	questionList := getSession(suite, w, tokenUser_Foo)
	question_FOO_Q := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Foo Question" })]

	updateQuestionDto := dtos.UpdateQuestionDto{Id: question_FOO_Q.Id, Text: "Updated Foo Question", Anonymous: true}

	assert.Equal(suite.T(), "Foo Question", question_FOO_Q.Text)
	assert.Equal(suite.T(), false, question_FOO_Q.Anonymous)

	putUpdateQuestion(suite, w, updateQuestionDto, tokenUser_Foo)

	updatedQuestionList := getSession(suite, w, tokenUser_Foo)
	updated_question_FOO_Q := updatedQuestionList[slices.IndexFunc(updatedQuestionList, func(c dtos.QuestionDto) bool { return c.Id == question_FOO_Q.Id })]

	assert.Equal(suite.T(), "Updated Foo Question", updated_question_FOO_Q.Text)
	assert.Equal(suite.T(), true, updated_question_FOO_Q.Anonymous)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *QuestionApiTestSuite) TestUpdateQuestion_EVENTS_RECEIVED_UPDATE_QUESTION() {
	w := httptest.NewRecorder()

	tokenUser_Foo := suite.tokenUser_Foo
	newQuestion := dtos.NewQuestionDto{Text: "Question to be updated", Anonymous: false}
	postNewQuestion(suite, w, newQuestion, tokenUser_Foo)
	questionList := getSession(suite, w, tokenUser_Foo)
	question_FOO_Q := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Question to be updated" })]

	updateQuestionDto := dtos.UpdateQuestionDto{Id: question_FOO_Q.Id, Text: "Updated Foo Question", Anonymous: true}

	putUpdateQuestion(suite, w, updateQuestionDto, tokenUser_Foo)

	assert.EventuallyWithT(suite.T(), func(c *assert.CollectT) {
		questionUpdatedEventFoo, err := findEvent[events.QuestionUpdated](suite.centrifugeClientFoo, events.UPDATE_QUESTION)
		assert.Nil(c, err)
		assert.Equal(c, "Updated Foo Question", questionUpdatedEventFoo.Text)
		assert.Equal(c, true, questionUpdatedEventFoo.Anonymous)
		assert.Equal(c, "", questionUpdatedEventFoo.Creator)

		questionUpdatedEventBar, err1 := findEvent[events.QuestionUpdated](suite.centrifugeClientBar, events.UPDATE_QUESTION)
		assert.Nil(c, err1)
		assert.Equal(c, "Updated Foo Question", questionUpdatedEventBar.Text)
		assert.Equal(c, true, questionUpdatedEventBar.Anonymous)
		assert.Equal(c, "", questionUpdatedEventBar.Creator)
	}, time.Second*3, time.Second*1, "Updated question event has not been received")

}

func (suite *QuestionApiTestSuite) TestUpdateQuestion_FORBIDDEN_403_WHEN_QUESTION_IS_NOT_OWNED() {
	w := httptest.NewRecorder()

	tokenUser_Foo := suite.tokenUser_Foo
	tokenUser_Bar := suite.tokenUser_Bar

	newQuestion := dtos.NewQuestionDto{Text: "Foo Question", Anonymous: false}
	postNewQuestion(suite, w, newQuestion, tokenUser_Foo)
	questionList := getSession(suite, w, tokenUser_Foo)
	question_FOO_Q := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Foo Question" })]

	updateQuestionDto := dtos.UpdateQuestionDto{Id: question_FOO_Q.Id, Text: "Updated Foo Question", Anonymous: true}

	assert.Equal(suite.T(), "Foo Question", question_FOO_Q.Text)
	assert.Equal(suite.T(), false, question_FOO_Q.Anonymous)

	w2 := httptest.NewRecorder()
	putUpdateQuestion(suite, w2, updateQuestionDto, tokenUser_Bar)

	assert.Equal(suite.T(), http.StatusForbidden, w2.Code)
}

func (suite *QuestionApiTestSuite) TestDeleteQuestion_OK_200_WHEN_QUESTION_IS_OWNED() {
	w := httptest.NewRecorder()

	tokenUser_Foo := suite.tokenUser_Foo
	newQuestion := dtos.NewQuestionDto{Text: "Foo Question", Anonymous: false}
	postNewQuestion(suite, w, newQuestion, tokenUser_Foo)
	questionList := getSession(suite, w, tokenUser_Foo)
	question_FOO_Q := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Foo Question" })]

	assert.Equal(suite.T(), "Foo Question", question_FOO_Q.Text)
	assert.Equal(suite.T(), false, question_FOO_Q.Anonymous)

	deleteQuestion(suite, w, question_FOO_Q.Id, tokenUser_Foo)

	updatedQuestionList := getSession(suite, w, tokenUser_Foo)
	idx := slices.IndexFunc(updatedQuestionList, func(c dtos.QuestionDto) bool { return c.Id == question_FOO_Q.Id })

	assert.Equal(suite.T(), -1, idx)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *QuestionApiTestSuite) TestDeleteQuestion_EVENTS_RECEIVED_DELETE_QUESTION() {
	w := httptest.NewRecorder()

	tokenUser_Foo := suite.tokenUser_Foo
	newQuestion := dtos.NewQuestionDto{Text: "Question to be deleted", Anonymous: false}
	postNewQuestion(suite, w, newQuestion, tokenUser_Foo)
	questionList := getSession(suite, w, tokenUser_Foo)
	question_FOO_Q := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Question to be deleted" })]

	deleteQuestion(suite, w, question_FOO_Q.Id, tokenUser_Foo)

	assert.EventuallyWithT(suite.T(), func(c *assert.CollectT) {
		questionUpdatedEventFoo, err := findEvent[events.QuestionDeleted](suite.centrifugeClientFoo, events.DELETE_QUESTION)
		assert.Nil(c, err)
		assert.Equal(c, question_FOO_Q.Id, questionUpdatedEventFoo.Id)

		questionUpdatedEventBar, err1 := findEvent[events.QuestionDeleted](suite.centrifugeClientBar, events.DELETE_QUESTION)
		assert.Nil(c, err1)
		assert.Equal(c, question_FOO_Q.Id, questionUpdatedEventBar.Id)
	}, time.Second*3, time.Second*1, "Deleted question event has not been received")
}

func (suite *QuestionApiTestSuite) TestDeleteQuestion_FORBIDDEN_403_WHEN_QUESTION_IS_NOT_OWNED() {
	w := httptest.NewRecorder()

	tokenUser_Foo := suite.tokenUser_Foo
	tokenUser_Bar := suite.tokenUser_Bar

	newQuestion := dtos.NewQuestionDto{Text: "Foo Question", Anonymous: false}
	postNewQuestion(suite, w, newQuestion, tokenUser_Foo)
	questionList := getSession(suite, w, tokenUser_Foo)
	question_FOO_Q := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Foo Question" })]

	assert.Equal(suite.T(), "Foo Question", question_FOO_Q.Text)
	assert.Equal(suite.T(), false, question_FOO_Q.Anonymous)

	w2 := httptest.NewRecorder()
	deleteQuestion(suite, w2, question_FOO_Q.Id, tokenUser_Bar)

	assert.Equal(suite.T(), http.StatusForbidden, w2.Code)
}

func (suite *QuestionApiTestSuite) TestAnswerQuestion_FORBIDDEN_403_WHEN_USER_NOT_SESSION_ADMIN_OR_ADMIN() {
	w := httptest.NewRecorder()

	tokenUser_Foo := suite.tokenUser_Foo
	tokenUser_Bar := suite.tokenUser_Bar

	newQuestion := dtos.NewQuestionDto{Text: "Foo Question", Anonymous: false}
	postNewQuestion(suite, w, newQuestion, tokenUser_Foo)
	questionList := getSession(suite, w, tokenUser_Foo)
	question_FOO_Q := questionList[slices.IndexFunc(questionList, func(c dtos.QuestionDto) bool { return c.Text == "Foo Question" })]

	assert.Equal(suite.T(), "Foo Question", question_FOO_Q.Text)
	assert.Equal(suite.T(), false, question_FOO_Q.Anonymous)

	w2 := httptest.NewRecorder()
	answerQuestion(suite, w2, question_FOO_Q.Id, tokenUser_Bar)

	assert.Equal(suite.T(), http.StatusForbidden, w2.Code)
}

func (suite *QuestionApiTestSuite) TestStartSession_FORBIDDEN_403_WHEN_USER_NOT_ADMIN() {
	type test struct {
		name  string
		token string
	}

	tests := []test{
		{"SESSION_ADMIN_403", suite.tokenUser_SessionAdmin},
		{"CONTRIBUTOR_403", suite.tokenUser_Bar},
	}

	for _, test := range tests {
		suite.T().Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", suite.apiPrefix, "/question/session/start"), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", test.token))
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusForbidden, w.Code)
		})
	}
}

func (suite *QuestionApiTestSuite) TestStopSession_FORBIDDEN_403_WHEN_USER_NOT_ADMIN() {
	type test struct {
		name  string
		token string
	}

	tests := []test{
		{"SESSION_ADMIN_403", suite.tokenUser_SessionAdmin},
		{"CONTRIBUTOR_403", suite.tokenUser_Bar},
	}

	for _, test := range tests {
		suite.T().Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", suite.apiPrefix, "/question/session/stop"), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", test.token))
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), http.StatusForbidden, w.Code)
		})
	}
}

func TestQuestionApiSuite(t *testing.T) {
	suite.Run(t, new(QuestionApiTestSuite))
}

func undoVoteQuestion(suite *QuestionApiTestSuite, w *httptest.ResponseRecorder, questionId, token string) {
	reqv, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/question/undovote/%s", suite.apiPrefix, questionId), nil)
	reqv.Header.Set("Content-Type", "application/json")
	reqv.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, reqv)
}

func upvoteQuestion(suite *QuestionApiTestSuite, w *httptest.ResponseRecorder, questionId, token string) {
	reqv, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/question/upvote/%s", suite.apiPrefix, questionId), nil)
	reqv.Header.Set("Content-Type", "application/json")
	reqv.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, reqv)
}

func answerQuestion(suite *QuestionApiTestSuite, w *httptest.ResponseRecorder, questionId, token string) {
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/question/answer/%s", suite.apiPrefix, questionId), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, req)
}

func postNewQuestion(suite *QuestionApiTestSuite, w *httptest.ResponseRecorder, question dtos.NewQuestionDto, token string) {
	newQuestion, _ := json.Marshal(question)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/question/new", suite.apiPrefix), bytes.NewBuffer(newQuestion))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, req)
}

func putUpdateQuestion(suite *QuestionApiTestSuite, w *httptest.ResponseRecorder, question dtos.UpdateQuestionDto, token string) {
	updateQuestion, _ := json.Marshal(question)
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/question/update", suite.apiPrefix), bytes.NewBuffer(updateQuestion))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, req)
}

func deleteQuestion(suite *QuestionApiTestSuite, w *httptest.ResponseRecorder, id string, token string) {
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/question/delete/%s", suite.apiPrefix, id), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, req)
}

func startSession(suite *QuestionApiTestSuite) {
	w := httptest.NewRecorder()

	tokenUser_Admin := suite.tokenUser_Admin

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/question/session/start", suite.apiPrefix), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenUser_Admin))
	suite.router.ServeHTTP(w, req)
}

func stopSession(suite *QuestionApiTestSuite) {
	w := httptest.NewRecorder()

	tokenUser_Admin := suite.tokenUser_Admin

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/question/session/stop", suite.apiPrefix), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenUser_Admin))
	suite.router.ServeHTTP(w, req)
}

func getSession(suite *QuestionApiTestSuite, w *httptest.ResponseRecorder, token string) []dtos.QuestionDto {
	reql, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/question/session", suite.apiPrefix), nil)
	reql.Header.Set("Content-Type", "application/json")
	reql.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	suite.router.ServeHTTP(w, reql)

	var questionList []dtos.QuestionDto
	body, _ := io.ReadAll(w.Body)
	json.Unmarshal(body, &questionList)

	return questionList
}

func findEvent[T any](centrifugeTestClient CentrifugeTestClient, eventType events.EventType) (*T, error) {
	for i := range centrifugeTestClient.receivedMessages {
		if centrifugeTestClient.receivedMessages[i].EventType == eventType {
			var eventPayload T
			json.Unmarshal([]byte(centrifugeTestClient.receivedMessages[i].Payload), &eventPayload)
			return &eventPayload, nil
		}
	}

	return nil, errors.New("Event not found")
}

func clearCentrifugeClientsMessages(suite *QuestionApiTestSuite) {
	suite.centrifugeClientFoo.receivedMessages = suite.centrifugeClientFoo.receivedMessages[:0]
	suite.centrifugeClientBar.receivedMessages = suite.centrifugeClientBar.receivedMessages[:0]
}

func initCentrifuge(suite *QuestionApiTestSuite) {
	time.Sleep(time.Second * 2)

	wsURL := "ws://localhost:3333/api/v1/connection/websocket"
	cFoo := centrifuge.NewJsonClient(wsURL, centrifuge.Config{
		Token: suite.tokenUser_Foo,
	})
	cBar := centrifuge.NewJsonClient(wsURL, centrifuge.Config{
		Token: suite.tokenUser_Bar,
	})

	errFoo := cFoo.Connect()
	errBar := cBar.Connect()

	if errFoo != nil {
		logrus.Fatal(errFoo)
	}

	if errBar != nil {
		logrus.Fatal(errBar)
	}

	suite.centrifugeClientFoo = CentrifugeTestClient{
		client:           cFoo,
		receivedMessages: make([]events.Event, 0),
	}
	suite.centrifugeClientBar = CentrifugeTestClient{
		client:           cBar,
		receivedMessages: make([]events.Event, 0),
	}
}
