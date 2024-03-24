package services

import (
	"net/http"
	"voting/dtos"
	"voting/internal/votingstorage"
	shared_infra_broker "voting/shared/infra/broker"
	"voting/shared/shared_models"

	"github.com/gin-gonic/gin"
)

type BrokeredQuestionsService struct {
	Broker          shared_infra_broker.Broker
	QuestionSession votingstorage.VotingStorage
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
