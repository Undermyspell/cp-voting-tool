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

// AddQuestion         godoc
// @Security 	 JWT
// @Summary      Deletes an existing question of the current session
// @Description  Deletes an existing question of the current session, only owned questions can be updated
// @Tags         Question
// @Produce      json
// @Param        id  path  string  true  "Id of question to delete"
// @Success      200
// @Failure      401
// @Failure      403 {string} error
// @Router       /api/v1/question/delete/{id} [delete]
func (service *BrokeredQuestionsService) Delete(c *gin.Context) {
	user, _ := c.Get(shared_models.User)
	userContext := user.(*shared_models.UserContext)
	questionId := c.Param("id")

	err := service.deleteQuestion(questionId, *userContext)

	if err != nil {
		c.JSON(int(err.HttpStatus), gin.H{
			"error": err.Error(),
		})
		return
	}

	questionDeletedSseMessage := events.QuestionDeleted{
		Id: questionId,
	}
	questionDeletedByteString, _ := json.Marshal(questionDeletedSseMessage)

	event := events.Event{
		EventType: events.DELETE_QUESTION,
		Payload:   string(questionDeletedByteString),
	}

	service.Broker.NotifyAll(event)
}

// UndovoteQuestion         godoc
// @Security 	 JWT
// @Summary      Undo the upvote for a question
// @Description  Undo the upvote of a question of the current session
// @Tags         Question
// @Produce      json
// @Param        id  path  string  true  "Id of question to undo the vote"
// @Success      200
// @Failure      401
// @Failure      404 {string} error
// @Router       /api/v1/question/undovote/{id} [put]
func (service *BrokeredQuestionsService) UndoVote(c *gin.Context) {
	user, _ := c.Get(shared_models.User)
	questionId := c.Param("id")
	userContext := user.(*shared_models.UserContext)

	votes, err := service.undoVote(questionId, *userContext)

	if err != nil {
		c.JSON(int(err.HttpStatus), gin.H{
			"error": err.Error(),
		})
		return
	}

	questionUndoUpvoteMessage := struct {
		Id    string
		Votes int
	}{questionId, votes}

	questionUndoUpVoteForUserMessage := events.QuestionUpvoted{
		Id:    questionId,
		Votes: votes,
		Voted: false,
	}

	questionForUserPaylod, errf := json.Marshal(questionUndoUpVoteForUserMessage)
	questionPayload, errj := json.Marshal(questionUndoUpvoteMessage)

	if errj != nil {
		c.JSON(http.StatusBadRequest, "cant marshal question")
		return
	}

	if errf != nil {
		c.JSON(http.StatusBadRequest, "cant marshal question")
		return
	}

	event := events.Event{
		EventType: events.UNDO_UPVOTE_QUESTION,
		Payload:   string(questionPayload),
	}

	userevent := events.Event{
		EventType: events.UNDO_UPVOTE_QUESTION,
		Payload:   string(questionForUserPaylod),
	}

	service.Broker.NotifyUser(userevent, *userContext)
	service.Broker.NotifyAllButUser(event, *userContext)
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

func (service *BrokeredQuestionsService) deleteQuestion(id string, creator shared_models.UserContext) *validation.ValidationError {
	if !service.QuestionSession.IsRunning() {
		return &validation.ValidationError{
			ValidationError: "no questions session currently running",
			HttpStatus:      http.StatusNotAcceptable,
		}
	}

	questionToDelete, ok := service.QuestionSession.GetQuestion(id)
	if !ok {
		return &validation.ValidationError{
			ValidationError: "question not found",
			HttpStatus:      http.StatusNotFound,
		}
	}

	if questionToDelete.CreatorHash != creator.GetHash(service.QuestionSession.GetSecret()) {
		return &validation.ValidationError{
			ValidationError: "you do not own this question",
			HttpStatus:      http.StatusForbidden,
		}
	}

	if questionToDelete.Answered {
		return &validation.ValidationError{
			ValidationError: "question has already been answered",
			HttpStatus:      http.StatusNotAcceptable,
		}
	}

	service.QuestionSession.DeleteQuestion(questionToDelete.Id)

	return nil
}

func (service *BrokeredQuestionsService) undoVote(id string, user shared_models.UserContext) (int, *validation.ValidationError) {
	if !service.QuestionSession.IsRunning() {
		return 0, &validation.ValidationError{
			ValidationError: "no questions session currently running",
			HttpStatus:      http.StatusNotAcceptable,
		}
	}

	question, ok := service.QuestionSession.GetQuestion(id)

	if !ok {
		return 0, &validation.ValidationError{
			ValidationError: "question not found",
			HttpStatus:      http.StatusNotFound,
		}
	}

	answered := question.Answered
	if answered {
		return 0, &validation.ValidationError{
			ValidationError: "question already answered",
			HttpStatus:      http.StatusNotAcceptable,
		}
	}

	hash := user.GetHash(service.QuestionSession.GetSecret())
	_, ok = service.QuestionSession.GetUserVotes()[hash][id]

	if !ok {
		return 0, &validation.ValidationError{
			ValidationError: "user has not voted",
			HttpStatus:      http.StatusNotAcceptable,
		}
	}

	service.QuestionSession.UndoVote(hash, id)

	question, _ = service.QuestionSession.GetQuestion(question.Id)

	return question.Votes, nil
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
