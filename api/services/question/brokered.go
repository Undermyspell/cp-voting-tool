package services

import (
	"encoding/json"
	"net/http"
	"sse/dtos"
	"sse/internal/broker"
	"sse/internal/helper"
	"sse/internal/models"
	"sse/internal/sse"
	"sse/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BrokeredQuestionsService struct {
	Broker        broker.Broker
	Session       map[string]*models.Question
	UserVotes     *models.SafeUserVotes
	SessionSecret string
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
// @Router       /api/v1/question/new [post]
func (service *BrokeredQuestionsService) Add(c *gin.Context) {
	var newQuestionDto dtos.NewQuestionDto
	user, _ := c.Get(models.User)

	userContext := user.(*models.UserContext)

	err := c.BindJSON(&newQuestionDto)

	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	question, errValidation := service.newQuestion(newQuestionDto.Text, newQuestionDto.Anonymous, *userContext)

	if errValidation != nil {
		c.JSON(int(errValidation.HttpStatus), gin.H{
			"error": errValidation.Error(),
		})
		return
	}

	newQuestionForUserSseMessage := dtos.QuestionDto{
		Id:        question.Id,
		Text:      question.Text,
		Creator:   question.CreatorName,
		Answered:  question.Answered,
		Votes:     question.Votes.Value(),
		Voted:     question.Voted,
		Anonymous: question.Anonymous,
		Owned:     true,
	}

	creatorForAllButUser := ""

	if !question.Anonymous {
		creatorForAllButUser = question.CreatorName
	}

	newQuestionForAllButUserSseMessage := dtos.QuestionDto{
		Id:        question.Id,
		Text:      question.Text,
		Creator:   creatorForAllButUser,
		Answered:  question.Answered,
		Votes:     question.Votes.Value(),
		Anonymous: question.Anonymous,
		Owned:     false,
	}

	newQuestionForUserByteString, _ := json.Marshal(newQuestionForUserSseMessage)
	newQuestionForAllButUserByteString, _ := json.Marshal(newQuestionForAllButUserSseMessage)

	eventForUser := sse.Event{
		EventType: sse.NEW_QUESTION,
		Payload:   string(newQuestionForUserByteString),
	}
	eventForAllButUser := sse.Event{
		EventType: sse.NEW_QUESTION,
		Payload:   string(newQuestionForAllButUserByteString),
	}

	service.Broker.NotifyUser(eventForUser, *userContext)
	service.Broker.NotifyAllButUser(eventForAllButUser, *userContext)
}

// UpdateQuestion         godoc
// @Security 	 JWT
// @Summary      Updates an existing question of the current session
// @Description  Updates an existing question of the current session, only owned questions can be updated
// @Tags         Question
// @Produce      json
// @Param        question  body      dtos.UpdateQuestionDto  true  "Question JSON"
// @Success      200
// @Failure      401
// @Failure      403
// @Router       /api/v1/question/update [put]
func (service *BrokeredQuestionsService) Update(c *gin.Context) {
	var updateQuestionDto dtos.UpdateQuestionDto
	user, _ := c.Get(models.User)
	userContext := user.(*models.UserContext)

	err := c.BindJSON(&updateQuestionDto)

	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	questionToUpdate, errValidation := service.updateQuestion(updateQuestionDto, *userContext)

	if errValidation != nil {
		c.JSON(int(errValidation.HttpStatus), gin.H{
			"error": errValidation.Error(),
		})
		return
	}

	questionToUpdateSseMessage := struct {
		Id        string
		Text      string
		Creator   string
		Anonymous bool
	}{questionToUpdate.Id, questionToUpdate.Text, questionToUpdate.CreatorName, questionToUpdate.Anonymous}

	newQuestionByteString, _ := json.Marshal(questionToUpdateSseMessage)

	event := sse.Event{
		EventType: sse.UPDATE_QUESTION,
		Payload:   string(newQuestionByteString),
	}

	service.Broker.NotifyAll(event)
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
	user, _ := c.Get(models.User)
	userContext := user.(*models.UserContext)
	questionId := c.Param("id")

	err := service.deleteQuestion(questionId, *userContext)

	if err != nil {
		c.JSON(int(err.HttpStatus), gin.H{
			"error": err.Error(),
		})
		return
	}

	questionDeletedSseMessage := struct {
		Id string
	}{questionId}
	questionDeletedByteString, _ := json.Marshal(questionDeletedSseMessage)

	event := sse.Event{
		EventType: sse.DELETE_QUESTION,
		Payload:   string(questionDeletedByteString),
	}

	service.Broker.NotifyAll(event)
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
// @Failure      404 {string} error
// @Router       /api/v1/question/upvote/{id} [put]
func (service *BrokeredQuestionsService) Upvote(c *gin.Context) {
	user, _ := c.Get(models.User)
	questionId := c.Param("id")
	userContext := user.(*models.UserContext)

	votes, err := service.upVote(questionId, *userContext)

	if err != nil {
		c.JSON(int(err.HttpStatus), gin.H{
			"error": err.Error(),
		})
		return
	}

	questionUpvoteSseMessage := struct {
		Id    string
		Votes int
	}{questionId, votes}

	questionUpVoteForUserSseMessage := struct {
		Id    string
		Votes int
		Voted bool
	}{questionId, votes, true}

	questionForUserPaylod, errf := json.Marshal(questionUpVoteForUserSseMessage)
	questionPayload, errj := json.Marshal(questionUpvoteSseMessage)

	if errj != nil {
		c.JSON(http.StatusBadRequest, "cant marshal question")
		return
	}

	if errf != nil {
		c.JSON(http.StatusBadRequest, "cant marshal question")
		return
	}

	event := sse.Event{
		EventType: sse.UPVOTE_QUESTION,
		Payload:   string(questionPayload),
	}

	userevent := sse.Event{
		EventType: sse.UPVOTE_QUESTION,
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

	questionMessage := struct {
		Id string
	}{questionId}
	questionPayload, errj := json.Marshal(questionMessage)

	if errj != nil {
		c.JSON(http.StatusBadRequest, "cant marshal question")
		return
	}

	event := sse.Event{
		EventType: sse.ANSWER_QUESTION,
		Payload:   string(questionPayload),
	}

	service.Broker.NotifyAll(event)
}

// StopSession         godoc
// @Security 	 JWT
// @Summary      Stops the current questions session
// @Description  Stops the current questions session
// @Tags         Question
// @Produce      json
// @Success      200
// @Failure      401
// @Router       /api/v1/question/session/stop [post]
func (service *BrokeredQuestionsService) Stop(c *gin.Context) {
	service.stop()

	event := sse.Event{
		EventType: sse.STOP_SESSION,
		Payload:   sse.PayloadEmpty,
	}

	service.Broker.NotifyAll(event)
}

// StartSession         godoc
// @Security 	 JWT
// @Summary      Starts a new questions session
// @Description  Starts a new questions session
// @Tags         Question
// @Produce      json
// @Success      200
// @Failure      401
// @Router       /api/v1/question/session/start [post]
func (service *BrokeredQuestionsService) Start(c *gin.Context) {
	service.start()

	event := sse.Event{
		EventType: sse.START_SESSION,
		Payload:   sse.PayloadEmpty,
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
	if service.Session == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "no questions session currently running",
		})
		return
	}

	user, _ := c.Get(models.User)
	userContext := user.(*models.UserContext)
	hash := userContext.GetHash(service.SessionSecret)
	questions := []dtos.QuestionDto{}

	for _, v := range service.Session {
		creator := v.CreatorName
		owned := v.CreatorHash == userContext.GetHash(service.SessionSecret)

		questions = append(questions, dtos.QuestionDto{
			Id:        v.Id,
			Text:      v.Text,
			Votes:     v.Votes.Value(),
			Answered:  v.Answered,
			Voted:     service.UserVotes.Value()[hash][v.Id],
			Anonymous: v.Anonymous,
			Creator:   creator,
			Owned:     owned,
		})
	}

	c.JSON(http.StatusOK, questions)
}

func (service *BrokeredQuestionsService) newQuestion(text string, anonymous bool, creator models.UserContext) (*models.Question, *validation.ValidationError) {
	if service.Session == nil {
		return nil, &validation.ValidationError{
			ValidationError: "no questions session currently running",
			HttpStatus:      http.StatusNotAcceptable,
		}
	}

	question := models.NewQuestion(text, anonymous, creator.Name, creator.GetHash(service.SessionSecret))
	service.Session[question.Id] = question
	return question, nil
}

func (service *BrokeredQuestionsService) updateQuestion(question dtos.UpdateQuestionDto, creator models.UserContext) (*models.Question, *validation.ValidationError) {
	if service.Session == nil {
		return nil, &validation.ValidationError{
			ValidationError: "no questions session currently running",
			HttpStatus:      http.StatusNotAcceptable,
		}
	}

	questionToUpdate, ok := service.Session[question.Id]
	if !ok {
		return nil, &validation.ValidationError{
			ValidationError: "question not found",
			HttpStatus:      http.StatusNotFound,
		}
	}

	if questionToUpdate.CreatorHash != creator.GetHash(service.SessionSecret) {
		return nil, &validation.ValidationError{
			ValidationError: "you do not own this question",
			HttpStatus:      http.StatusForbidden,
		}
	}

	questionToUpdate.Text = question.Text
	questionToUpdate.Anonymous = question.Anonymous
	questionToUpdate.CreatorName = ""
	if !questionToUpdate.Anonymous {
		questionToUpdate.CreatorName = creator.Name
	}

	return questionToUpdate, nil
}

func (service *BrokeredQuestionsService) deleteQuestion(id string, creator models.UserContext) *validation.ValidationError {
	if service.Session == nil {
		return &validation.ValidationError{
			ValidationError: "no questions session currently running",
			HttpStatus:      http.StatusNotAcceptable,
		}
	}

	questionToDelete, ok := service.Session[id]
	if !ok {
		return &validation.ValidationError{
			ValidationError: "question not found",
			HttpStatus:      http.StatusNotFound,
		}
	}

	if questionToDelete.CreatorHash != creator.GetHash(service.SessionSecret) {
		return &validation.ValidationError{
			ValidationError: "you do not own this question",
			HttpStatus:      http.StatusForbidden,
		}
	}

	delete(service.Session, questionToDelete.Id)

	return nil
}

func (service *BrokeredQuestionsService) upVote(id string, user models.UserContext) (int, *validation.ValidationError) {
	if service.Session == nil {
		return 0, &validation.ValidationError{
			ValidationError: "no questions session currently running",
			HttpStatus:      http.StatusNotAcceptable,
		}
	}

	question, ok := service.Session[id]

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

	hash := user.GetHash(service.SessionSecret)
	_, ok = service.UserVotes.Value()[hash][id]

	if ok {
		return 0, &validation.ValidationError{
			ValidationError: "user already voted",
			HttpStatus:      http.StatusNotAcceptable,
		}
	}

	service.UserVotes.SetUserVote(hash, id)
	question.Votes.Increment()

	return question.Votes.Value(), nil
}

func (service *BrokeredQuestionsService) answer(id string) *validation.ValidationError {
	if service.Session == nil {
		return &validation.ValidationError{
			ValidationError: "no questions session currently running",
			HttpStatus:      http.StatusNotAcceptable,
		}
	}

	question, ok := service.Session[id]

	if !ok {
		return &validation.ValidationError{
			ValidationError: "question not found",
			HttpStatus:      http.StatusNotFound,
		}
	}

	question.Answered = true

	return nil
}

func (service *BrokeredQuestionsService) stop() {
	service.UserVotes = nil
	service.Session = nil
	service.SessionSecret = ""
}

func (service *BrokeredQuestionsService) start() {
	service.UserVotes = models.NewSafeUserVotes()
	service.Session = make(map[string]*models.Question)
	service.SessionSecret = helper.GetRandomId()
}
