package votinghttp

import (
	"net/http"
	"voting/internal/models"
	usecases "voting/voting/use-cases/upvote-question"

	"github.com/gin-gonic/gin"
)

func Upvote(c *gin.Context) {
	user, _ := c.Get(models.User)
	questionId := c.Param("id")
	userContext := user.(*models.UserContext)

	err := usecases.Upvote(questionId, *userContext)

	httpStatus := http.StatusOK
	if err != nil {
		switch err.(type) {
		case *usecases.QuestionNotFoundError:
			httpStatus = http.StatusNotFound
		case *usecases.QuestionAlreadyAnsweredError:
			httpStatus = http.StatusNotAcceptable
		case *usecases.UserAlreadyVotedError:
			httpStatus = http.StatusNotAcceptable
		case *usecases.QuestionSessionNotRunningError:
			httpStatus = http.StatusNotAcceptable
		case *usecases.UnexpectedError:
			httpStatus = http.StatusBadRequest
		}
	}

	if err != nil {
		c.JSON(int(httpStatus), gin.H{
			"error": err.Error(),
		})
	}
}
