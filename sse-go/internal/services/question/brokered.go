package services

import (
	"encoding/json"
	"log"
	"net/http"
	"sse/internal/broker"
	"sse/internal/models"
	"sse/internal/sse"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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

	question := models.New(message.Text)

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

func (service *BrokeredQuestionsService) GetToken(c *gin.Context) {
	sampleSecretKey := []byte("my_test_secret")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().UTC().Add(time.Second * 3600).Unix()
	claims["user"] = "Hoodini Magician"
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		log.Println("Signing error")
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
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
