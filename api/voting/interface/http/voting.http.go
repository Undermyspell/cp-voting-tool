package votinghttp

import (
	"net/http"
	"voting/dtos"
	"voting/shared/shared_models"
	usecaseErrors "voting/voting/use-cases/_errors"
	ucAnswer "voting/voting/use-cases/answer-question"
	ucCreate "voting/voting/use-cases/create-question"
	ucDelete "voting/voting/use-cases/delete-question"
	ucGetSession "voting/voting/use-cases/get-session"
	ucStart "voting/voting/use-cases/start-session"
	ucStop "voting/voting/use-cases/stop-session"
	ucUndoVote "voting/voting/use-cases/undovote-question"
	ucUpdate "voting/voting/use-cases/update-question"
	ucUpvote "voting/voting/use-cases/upvote-question"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetSession(c *gin.Context) {
	user, _ := c.Get(shared_models.User)
	userContext := user.(*shared_models.UserContext)

	questions, err := ucGetSession.GetSession(userContext)

	httpStatus := http.StatusOK
	if err != nil {
		switch err.(type) {
		case *usecaseErrors.QuestionSessionNotRunningError:
			httpStatus = http.StatusNotAcceptable
		case *usecaseErrors.UnexpectedError:
			httpStatus = http.StatusBadRequest
		}

		c.JSON(int(httpStatus), gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, questions)
}

func StartSession(c *gin.Context) {
	ucStart.StartSession()
}

func StopSession(c *gin.Context) {
	ucStop.StopSession()
}

func Create(c *gin.Context) {
	var newQuestionDto dtos.NewQuestionDto
	user, _ := c.Get(shared_models.User)

	userContext := user.(*shared_models.UserContext)

	err := c.BindJSON(&newQuestionDto)

	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	err = ucCreate.Create(newQuestionDto, *userContext)

	httpStatus := http.StatusOK
	if err != nil {
		switch err.(type) {
		case *usecaseErrors.QuestionMaxLengthExceededError:
			httpStatus = http.StatusBadRequest
		case *usecaseErrors.QuestionSessionNotRunningError:
			httpStatus = http.StatusNotAcceptable
		case *usecaseErrors.UnexpectedError:
			httpStatus = http.StatusBadRequest
		}
	}

	if err != nil {
		c.JSON(int(httpStatus), gin.H{
			"error": err.Error(),
		})
	}
}

func Update(c *gin.Context) {
	var updateQuestionDto dtos.UpdateQuestionDto
	user, _ := c.Get(shared_models.User)
	userContext := user.(*shared_models.UserContext)

	err := c.BindJSON(&updateQuestionDto)

	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	err = ucUpdate.UpdateQuestion(updateQuestionDto, *userContext)

	httpStatus := http.StatusOK
	if err != nil {
		switch err.(type) {
		case *usecaseErrors.QuestionNotFoundError:
			httpStatus = http.StatusNotFound
		case *usecaseErrors.QuestionAlreadyAnsweredError:
			httpStatus = http.StatusNotAcceptable
		case *usecaseErrors.QuestionSessionNotRunningError:
			httpStatus = http.StatusNotAcceptable
		case *usecaseErrors.QuestionNotOwnedError:
			httpStatus = http.StatusForbidden
		case *usecaseErrors.UnexpectedError:
			httpStatus = http.StatusBadRequest
		}
	}

	if err != nil {
		c.JSON(int(httpStatus), gin.H{
			"error": err.Error(),
		})
	}
}

func Delete(c *gin.Context) {
	user, _ := c.Get(shared_models.User)
	userContext := user.(*shared_models.UserContext)
	questionId := c.Param("id")

	err := ucDelete.Delete(questionId, *userContext)

	httpStatus := http.StatusOK
	if err != nil {
		switch err.(type) {
		case *usecaseErrors.QuestionNotFoundError:
			httpStatus = http.StatusNotFound
		case *usecaseErrors.QuestionAlreadyAnsweredError:
			httpStatus = http.StatusNotAcceptable
		case *usecaseErrors.QuestionSessionNotRunningError:
			httpStatus = http.StatusNotAcceptable
		case *usecaseErrors.QuestionNotOwnedError:
			httpStatus = http.StatusForbidden
		case *usecaseErrors.UnexpectedError:
			httpStatus = http.StatusBadRequest
		}
	}

	if err != nil {
		c.JSON(int(httpStatus), gin.H{
			"error": err.Error(),
		})
	}
}

func Upvote(c *gin.Context) {
	user, _ := c.Get(shared_models.User)
	questionId := c.Param("id")
	userContext := user.(*shared_models.UserContext)

	err := ucUpvote.Upvote(questionId, *userContext)

	httpStatus := http.StatusOK
	if err != nil {
		switch err.(type) {
		case *usecaseErrors.QuestionNotFoundError:
			httpStatus = http.StatusNotFound
		case *usecaseErrors.QuestionAlreadyAnsweredError:
			httpStatus = http.StatusNotAcceptable
		case *usecaseErrors.UserAlreadyVotedError:
			httpStatus = http.StatusNotAcceptable
		case *usecaseErrors.QuestionSessionNotRunningError:
			httpStatus = http.StatusNotAcceptable
		case *usecaseErrors.UnexpectedError:
			httpStatus = http.StatusBadRequest
		}
	}

	if err != nil {
		c.JSON(int(httpStatus), gin.H{
			"error": err.Error(),
		})
	}
}

func UndoVote(c *gin.Context) {
	user, _ := c.Get(shared_models.User)
	questionId := c.Param("id")
	userContext := user.(*shared_models.UserContext)

	err := ucUndoVote.UndoVote(questionId, *userContext)

	httpStatus := http.StatusOK
	if err != nil {
		switch err.(type) {
		case *usecaseErrors.QuestionNotFoundError:
			httpStatus = http.StatusNotFound
		case *usecaseErrors.QuestionAlreadyAnsweredError:
			httpStatus = http.StatusNotAcceptable
		case *usecaseErrors.UserHasNotVotedError:
			httpStatus = http.StatusNotAcceptable
		case *usecaseErrors.QuestionSessionNotRunningError:
			httpStatus = http.StatusNotAcceptable
		case *usecaseErrors.UnexpectedError:
			httpStatus = http.StatusBadRequest
		}
	}

	if err != nil {
		c.JSON(int(httpStatus), gin.H{
			"error": err.Error(),
		})
	}
}

func Answer(c *gin.Context) {
	questionId := c.Param("id")

	err := ucAnswer.Answer(questionId)

	httpStatus := http.StatusOK
	if err != nil {
		switch err.(type) {
		case *usecaseErrors.QuestionNotFoundError:
			httpStatus = http.StatusNotFound
		case *usecaseErrors.QuestionAlreadyAnsweredError:
			httpStatus = http.StatusNotAcceptable
		case *usecaseErrors.QuestionSessionNotRunningError:
			httpStatus = http.StatusNotAcceptable
		case *usecaseErrors.UnexpectedError:
			httpStatus = http.StatusBadRequest
		}
	}

	if err != nil {
		c.JSON(int(httpStatus), gin.H{
			"error": err.Error(),
		})
	}
}
