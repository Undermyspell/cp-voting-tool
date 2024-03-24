package services

import (
	"encoding/json"
	"net/http"
	"voting/dtos"
	"voting/internal/events"
	"voting/internal/validation"
	"voting/internal/votingstorage"
	shared_infra_broker "voting/shared/infra/broker"
	"voting/shared/shared_models"

	"github.com/gin-gonic/gin"
)

type BrokeredQuestionsService struct {
	Broker          shared_infra_broker.Broker
	QuestionSession votingstorage.VotingStorage
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
// @Failure      404 {string} error
// @Router       /api/v1/question/answer/{id} [put]
func (service *BrokeredQuestionsService) Answer(c *gin.Context) {
	questionId := c.Param("id")

	err := service.answer(questionId)

	if err != nil {
		c.JSON(int(err.HttpStatus), gin.H{
			"error": err.Error(),
		})
		return
	}

	questionMessage := events.QuestionAnswered{
		Id: questionId,
	}
	questionPayload, errj := json.Marshal(questionMessage)

	if errj != nil {
		c.JSON(http.StatusBadRequest, "cant marshal question")
		return
	}

	event := events.Event{
		EventType: events.ANSWER_QUESTION,
		Payload:   string(questionPayload),
	}

	service.Broker.NotifyAll(event)
}

// GetSession         godoc
// @Security 	 JWT
// @Summary      Gets the questions of the current session
// @Description  Gets the questions of the current session
// @Tags         Question
// @Produce      json
// @Success      200 {array} dtos.QuestionDto
// @Failure      401
// @Router       /api/v1/question/session [get]
func (service *BrokeredQuestionsService) GetSession(c *gin.Context) {
	if !service.QuestionSession.IsRunning() {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "no questions session currently running",
		})
		return
	}

	user, _ := c.Get(shared_models.User)
	userContext := user.(*shared_models.UserContext)
	hash := userContext.GetHash(service.QuestionSession.GetSecret())
	questions := []dtos.QuestionDto{}

	for _, v := range service.QuestionSession.GetQuestions() {
		creator := v.CreatorName
		owned := v.CreatorHash == userContext.GetHash(service.QuestionSession.GetSecret())

		questions = append(questions, dtos.QuestionDto{
			Id:        v.Id,
			Text:      v.Text,
			Votes:     v.Votes,
			Answered:  v.Answered,
			Voted:     service.QuestionSession.GetUserVotes()[hash][v.Id],
			Anonymous: v.Anonymous,
			Creator:   creator,
			Owned:     owned,
		})
	}

	c.JSON(http.StatusOK, questions)
}

func (service *BrokeredQuestionsService) answer(id string) *validation.ValidationError {
	if !service.QuestionSession.IsRunning() {
		return &validation.ValidationError{
			ValidationError: "no questions session currently running",
			HttpStatus:      http.StatusNotAcceptable,
		}
	}

	_, ok := service.QuestionSession.GetQuestion(id)

	if !ok {
		return &validation.ValidationError{
			ValidationError: "question not found",
			HttpStatus:      http.StatusNotFound,
		}
	}

	service.QuestionSession.AnswerQuestion(id)

	return nil
}
