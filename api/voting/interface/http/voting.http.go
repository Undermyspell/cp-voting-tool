package votinghttp

import (
	"net/http"
	"voting/dtos"
	"voting/shared/shared_models"
	usecaseErrors "voting/voting/use-cases/_errors"
	ucCreate "voting/voting/use-cases/create-question"
	ucStart "voting/voting/use-cases/start-session"
	ucStop "voting/voting/use-cases/stop-session"
	ucUpdate "voting/voting/use-cases/update-question"
	ucUpvote "voting/voting/use-cases/upvote-question"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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
