package votinghttp

import (
	"net/http"
	"voting/dtos"
	"voting/internal/models"
	usecaseErrors "voting/voting/use-cases/_errors"
	ucCreate "voting/voting/use-cases/create-question"
	usecases "voting/voting/use-cases/upvote-question"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Create(c *gin.Context) {
	var newQuestionDto dtos.NewQuestionDto
	user, _ := c.Get(models.User)

	userContext := user.(*models.UserContext)

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

func Upvote(c *gin.Context) {
	user, _ := c.Get(models.User)
	questionId := c.Param("id")
	userContext := user.(*models.UserContext)

	err := usecases.Upvote(questionId, *userContext)

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
