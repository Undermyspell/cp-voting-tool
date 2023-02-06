package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"sse/dtos"
	"sse/internal/broker"
	"sse/internal/sse"
	"sse/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BrokeredQuestionsService struct {
	Broker  broker.Broker
	Session map[string]models.Question
}

// AddQuestion         godoc
// @Security 	 JWT
// @Summary      Adds a new question
// @Description  Adds a new question to the current session
// @Tags         Question
// @Produce      json
// @Param        question  body      dtos.NewQuestionDto  true  "Question JSON"
// @Success      200
// @Failure      401
// @Router       /question/new [post]
func (service *BrokeredQuestionsService) AddQuestion(c *gin.Context) {
	var message dtos.NewQuestionDto

	err := c.BindJSON(&message)

	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	question := models.New(message.Text)

	newQuestion, _ := json.Marshal(question)

	event := sse.Event{
		EventType: sse.NEW_QUESTION,
		Payload:   string(newQuestion),
	}

	service.Session[question.Id] = question

	service.Broker.Notify(event)
}

// UpvoteQuestion         godoc
// @Security 	 JWT
// @Summary      Upvotes a question
// @Description  Upvotes a question of the current session
// @Tags         Question
// @Produce      json
// @Param        id  path  string  true  "Id of question to upvote"
// @Success      200
// @Failure      401
// @Router       /question/upvote/{id} [put]
func (service *BrokeredQuestionsService) UpvoteQuestion(c *gin.Context) {
	questionId := c.Param("id")

	votes, err := service.upVote(questionId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	questionUpvoteSeeMessage := struct {
		Id    string
		Votes int
	}{questionId, votes}
	questionPayload, err := json.Marshal(questionUpvoteSeeMessage)

	if err != nil {
		c.JSON(http.StatusBadRequest, "cant marshal question")
		return
	}

	event := sse.Event{
		EventType: sse.UPVOTE_QUESTION,
		Payload:   string(questionPayload),
	}

	service.Broker.Notify(event)
}

// AnswerQuestion         godoc
// @Security 	 JWT
// @Summary      Answers a question
// @Description  Answers a question of the current session
// @Tags         Question
// @Produce      json
// @Param        id  path  string  true  "Id of question to answer"
// @Success      200
// @Failure      401
// @Router       /question/answer/{id} [put]
func (service *BrokeredQuestionsService) Answer(c *gin.Context) {
	questionId := c.Param("id")

	err := service.answer(questionId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	questionMessage := struct {
		Id string
	}{questionId}
	questionPayload, err := json.Marshal(questionMessage)

	if err != nil {
		c.JSON(http.StatusBadRequest, "cant marshal question")
		return
	}

	event := sse.Event{
		EventType: sse.ANSWER_QUESTION,
		Payload:   string(questionPayload),
	}

	service.Broker.Notify(event)
}

// ResetSession         godoc
// @Security 	 JWT
// @Summary      Resets the current session
// @Description  Resets the current question session
// @Tags         Question
// @Produce      json
// @Success      200
// @Failure      401
// @Router       /question/reset/ [post]
func (service *BrokeredQuestionsService) Reset(c *gin.Context) {
	service.reset()

	event := sse.Event{
		EventType: sse.ANSWER_QUESTION,
		Payload:   sse.PayloadEmpty,
	}

	service.Broker.Notify(event)
}

func (service *BrokeredQuestionsService) upVote(id string) (int, error) {
	question, ok := service.Session[id]
	if !ok {
		return 0, errors.New("question not found")
	}

	question.Votes++
	service.Session[id] = question

	return question.Votes, nil
}

func (service *BrokeredQuestionsService) answer(id string) error {
	question, ok := service.Session[id]

	if !ok {
		return errors.New("question not found")
	}

	question.Answered = true
	service.Session[id] = question

	return nil
}

func (service *BrokeredQuestionsService) reset() {
	service.Session = make(map[string]models.Question)
}
