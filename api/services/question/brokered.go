package services

import (
	"encoding/json"
	"net/http"
	"sse/dtos"
	"sse/internal/broker"
	"sse/internal/models"
	"sse/internal/sse"
	"sse/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BrokeredQuestionsService struct {
	Broker    broker.Broker
	Session   map[string]models.Question
	UserVotes map[string]map[string]bool
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

	newQuestionByteString, _ := json.Marshal(*question)

	event := sse.Event{
		EventType: sse.NEW_QUESTION,
		Payload:   string(newQuestionByteString),
	}

	service.Broker.Notify(event)
}

// AddQuestion         godoc
// @Security 	 JWT
// @Summary      Updates an existing question of the current session
// @Description  Updates an existing question of the current session, only owned questions can be updated
// @Tags         Question
// @Produce      json
// @Param        question  body      dtos.UpdateQuestionDto  true  "Question JSON"
// @Success      200
// @Failure      401
// @Failure      403
// @Router       /question/update [put]
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

	newQuestionByteString, _ := json.Marshal(*questionToUpdate)

	event := sse.Event{
		EventType: sse.UPDATE_QUESTION,
		Payload:   string(newQuestionByteString),
	}

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
// @Failure      404 {string} error
// @Router       /question/upvote/{id} [put]
func (service *BrokeredQuestionsService) Upvote(c *gin.Context) {
	user, _ := c.Get(models.User)
	questionId := c.Param("id")

	votes, err := service.upVote(user.(*models.UserContext), questionId)

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

	questionPayload, errj := json.Marshal(questionUpvoteSseMessage)

	if errj != nil {
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
// @Failure      404 {string} error
// @Router       /question/answer/{id} [put]
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

	service.Broker.Notify(event)
}

// StopSession         godoc
// @Security 	 JWT
// @Summary      Stops the current questions session
// @Description  Stops the current questions session
// @Tags         Question
// @Produce      json
// @Success      200
// @Failure      401
// @Router       /question/session/stop [post]
func (service *BrokeredQuestionsService) Stop(c *gin.Context) {
	service.stop()

	event := sse.Event{
		EventType: sse.STOP_SESSION,
		Payload:   sse.PayloadEmpty,
	}

	service.Broker.Notify(event)
}

// StartSession         godoc
// @Security 	 JWT
// @Summary      Starts a new questions session
// @Description  Starts a new questions session
// @Tags         Question
// @Produce      json
// @Success      200
// @Failure      401
// @Router       /question/session/start [post]
func (service *BrokeredQuestionsService) Start(c *gin.Context) {
	service.start()

	event := sse.Event{
		EventType: sse.START_SESSION,
		Payload:   sse.PayloadEmpty,
	}

	service.Broker.Notify(event)
}

// GetSession         godoc
// @Security 	 JWT
// @Summary      Gets the questions of the current session
// @Description  Gets the questions of the current session
// @Tags         Question
// @Produce      json
// @Success      200 {array} dtos.QuestionDto
// @Failure      401
// @Router       /question/session [get]
func (service *BrokeredQuestionsService) GetSession(c *gin.Context) {
	if service.Session == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "no questions session currently running",
		})
		return
	}

	user, _ := c.Get(models.User)
	userContext := user.(*models.UserContext)

	questions := []dtos.QuestionDto{}

	for _, v := range service.Session {
		creator := v.Creator.Name
		owned := v.Creator.Email == userContext.Email
		if v.Anonymous && !owned {
			creator = ""
		}

		questions = append(questions, dtos.QuestionDto{
			Id:        v.Id,
			Text:      v.Text,
			Votes:     v.Votes,
			Answered:  v.Answered,
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

	question := models.NewQuestion(text, anonymous, creator)
	service.Session[question.Id] = question
	return &question, nil
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

	if questionToUpdate.Creator.Email != creator.Email {
		return nil, &validation.ValidationError{
			ValidationError: "you do not own this question",
			HttpStatus:      http.StatusForbidden,
		}
	}

	questionToUpdate.Text = question.Text
	questionToUpdate.Anonymous = question.Anonymous

	service.Session[questionToUpdate.Id] = questionToUpdate
	return &questionToUpdate, nil
}

func (service *BrokeredQuestionsService) upVote(user *models.UserContext, id string) (int, *validation.ValidationError) {
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

	hash := user.GetHash()
	_, ok = service.UserVotes[hash][id]

	if ok {
		return 0, &validation.ValidationError{
			ValidationError: "user already voted",
			HttpStatus:      http.StatusNotAcceptable,
		}
	}

	_, ok = service.UserVotes[hash]

	if !ok {
		service.UserVotes[hash] = make(map[string]bool)
	}

	service.UserVotes[hash][id] = true

	question.Votes++
	service.Session[id] = question

	return question.Votes, nil
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
	service.Session[id] = question

	return nil
}

func (service *BrokeredQuestionsService) stop() {
	service.UserVotes = nil
	service.Session = nil
}

func (service *BrokeredQuestionsService) start() {
	service.UserVotes = make(map[string]map[string]bool)
	service.Session = make(map[string]models.Question)
}
