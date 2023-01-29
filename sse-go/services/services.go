package services

import (
	"encoding/json"
	"log"
	"net/http"
	"sse/internal/broker"
	"sse/internal/sse"
	"sse/models"

	"github.com/gin-gonic/gin"
)

type QuestionsService struct {
	Broker  *broker.Broker
	Session map[string]models.Question
}

func New(broker *broker.Broker) *QuestionsService {
	return &QuestionsService{
		Broker:  broker,
		Session: make(map[string]models.Question),
	}
}

func (service *QuestionsService) AddQuestion(c *gin.Context) {
	var message models.Question
	err := c.BindJSON(&message)
	newQuestion, _ := json.Marshal(message)

	if err != nil {
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	event := sse.Event{
		EventType: sse.NEW_QUESTION,
		Payload:   string(newQuestion),
	}

	service.Session[message.Id] = message

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

func (service *QuestionsService) upVote(id string) int {
	question := service.Session[id]
	question.Votes++
	service.Session[id] = question

	return question.Votes
}
