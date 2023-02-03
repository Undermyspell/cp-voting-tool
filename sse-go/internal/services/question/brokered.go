package services

import (
	"encoding/json"
	"log"
	"net/http"
	"sse/internal/broker"
	"sse/internal/models"
	"sse/internal/sse"

	"github.com/gin-gonic/gin"
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
		log.Println(err)
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

	log.Default().Println(service.Session)

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

	service.Broker.Notify(event)
}

func (service *BrokeredQuestionsService) Reset(c *gin.Context) {
	service.reset()
	log.Default().Print(service.Session)

	event := sse.Event{
		EventType: sse.ANSWER_QUESTION,
		Payload:   sse.PayloadEmpty,
	}

	service.Broker.Notify(event)
}

func (service *BrokeredQuestionsService) upVote(id string) int {
	question := service.Session[id]
	question.Votes++
	service.Session[id] = question

	return question.Votes
}

func (service *BrokeredQuestionsService) answer(id string) {
	question := service.Session[id]
	question.Answered = true
	service.Session[id] = question
}

func (service *BrokeredQuestionsService) reset() {
	service.Session = make(map[string]models.Question)
}
