package services

import (
	"encoding/json"
	"log"
	"net/http"
	"sse/internal/broker"
	"sse/internal/sse"
	"sse/models/question"

	"github.com/gin-gonic/gin"
)

type QuestionsService struct {
	Broker  *broker.Broker
	Session map[string]question.Question
}

func New(broker *broker.Broker) *QuestionsService {
	return &QuestionsService{
		Broker:  broker,
		Session: make(map[string]question.Question),
	}
}

func (service *QuestionsService) AddQuestion(c *gin.Context) {
	var message struct {
		Text string
	}
	err := c.BindJSON(&message)

	question := question.New(message.Text)

	newQuestion, _ := json.Marshal(question)

	if err != nil {
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	event := sse.Event{
		EventType: sse.NEW_QUESTION,
		Payload:   string(newQuestion),
	}

	service.Session[question.Id] = question

	log.Default().Println(service.Session)

	service.Broker.Notifier <- event
}

func (service *QuestionsService) UpvoteQuestion(c *gin.Context) {
	var questionMessage struct {
		Id    string
		Votes int
	}
	err := c.BindJSON(&questionMessage)

	if err != nil {
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	questionMessage.Votes = service.upVote(questionMessage.Id)
	questionPayload, err := json.Marshal(questionMessage)

	if err != nil {
		c.JSON(http.StatusBadRequest, "cant marshal question")
		return
	}

	event := sse.Event{
		EventType: sse.UPVOTE_QUESTION,
		Payload:   string(questionPayload),
	}

	service.Broker.Notifier <- event
}

func (service *QuestionsService) Answer(c *gin.Context) {
	var questionMessage struct {
		Id string
	}
	err := c.BindJSON(&questionMessage)

	if err != nil {
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	service.answer(questionMessage.Id)
	questionPayload, err := json.Marshal(questionMessage)

	log.Default().Print(service.Session)

	if err != nil {
		c.JSON(http.StatusBadRequest, "cant marshal question")
		return
	}

	event := sse.Event{
		EventType: sse.ANSWER_QUESTION,
		Payload:   string(questionPayload),
	}

	service.Broker.Notifier <- event
}

func (service *QuestionsService) Reset(c *gin.Context) {
	service.reset()
	log.Default().Print(service.Session)

	event := sse.Event{
		EventType: sse.ANSWER_QUESTION,
		Payload:   sse.PayloadEmpty,
	}

	service.Broker.Notifier <- event
}

func (service *QuestionsService) upVote(id string) int {
	question := service.Session[id]
	question.Votes++
	service.Session[id] = question

	return question.Votes
}

func (service *QuestionsService) answer(id string) {
	question := service.Session[id]
	question.Answered = true
	service.Session[id] = question
}

func (service *QuestionsService) reset() {
	service.Session = make(map[string]question.Question)
}
