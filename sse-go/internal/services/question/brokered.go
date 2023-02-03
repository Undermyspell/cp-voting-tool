package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"sse/internal/broker"
	"sse/internal/models"
	"sse/internal/sse"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BrokeredQuestionsService struct {
	Broker  broker.Broker
	Session map[string]models.Question
}

func (service *BrokeredQuestionsService) AddQuestion(c *gin.Context) {
	var message struct {
		Text string
	}
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

func (service *BrokeredQuestionsService) UpvoteQuestion(c *gin.Context) {
	var questionMessage struct {
		Id    string
		Votes int
	}
	err := c.BindJSON(&questionMessage)

	if err != nil {
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	questionMessage.Votes, err = service.upVote(questionMessage.Id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	questionPayload, err := json.Marshal(questionMessage)

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

func (service *BrokeredQuestionsService) Answer(c *gin.Context) {
	var questionMessage struct {
		Id string
	}
	err := c.BindJSON(&questionMessage)

	if err != nil {
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	err = service.answer(questionMessage.Id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

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
