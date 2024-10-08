package voting_http

import (
	"errors"
	"net/http"
	shared_models "voting/shared/models"
	usecases "voting/voting/usecases"
	usecaseErrors "voting/voting/usecases/_errors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetSession(c *gin.Context) {
	user, _ := c.Get(shared_models.User)
	userContext := user.(*shared_models.UserContext)

	questions, err := usecases.GetSession(userContext)

	httpStatus := http.StatusOK
	if err != nil {
		if errors.Is(err, usecaseErrors.ErrQuestionSessionNotRunning) {
			httpStatus = http.StatusNotAcceptable
		}
		if errors.Is(err, usecaseErrors.ErrUnexpected) {
			httpStatus = http.StatusBadRequest
		}

		c.JSON(int(httpStatus), gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, questions)
}

func StartSession(c *gin.Context) {
	usecases.StartSession()
}

func StopSession(c *gin.Context) {
	usecases.StopSession()
}

func Create(c *gin.Context) {
	var newQuestionDto usecases.NewQuestionDto
	user, _ := c.Get(shared_models.User)

	userContext := user.(*shared_models.UserContext)

	err := c.BindJSON(&newQuestionDto)

	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	err = usecases.Create(newQuestionDto, *userContext)

	httpStatus := http.StatusOK
	if err != nil {
		if errors.Is(err, usecaseErrors.ErrQuestionMaxLengthExceeded) {
			httpStatus = http.StatusBadRequest
		}
		if errors.Is(err, usecaseErrors.ErrQuestionSessionNotRunning) {
			httpStatus = http.StatusNotAcceptable
		}
		if errors.Is(err, usecaseErrors.ErrUnexpected) {
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
	var updateQuestionDto usecases.UpdateQuestionDto
	user, _ := c.Get(shared_models.User)
	userContext := user.(*shared_models.UserContext)

	err := c.BindJSON(&updateQuestionDto)

	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	err = usecases.UpdateQuestion(updateQuestionDto, *userContext)

	httpStatus := http.StatusOK
	if err != nil {
		if errors.Is(err, usecaseErrors.ErrQuestionNotFound) {
			httpStatus = http.StatusNotFound
		}
		if errors.Is(err, usecaseErrors.ErrQuestionAlreadyAnswered) {
			httpStatus = http.StatusNotAcceptable
		}
		if errors.Is(err, usecaseErrors.ErrQuestionSessionNotRunning) {
			httpStatus = http.StatusNotAcceptable
		}
		if errors.Is(err, usecaseErrors.ErrQuestionNotOwned) {
			httpStatus = http.StatusForbidden
		}
		if errors.Is(err, usecaseErrors.ErrUnexpected) {
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

	err := usecases.Delete(questionId, *userContext)

	httpStatus := http.StatusOK
	if err != nil {
		if errors.Is(err, usecaseErrors.ErrQuestionNotFound) {
			httpStatus = http.StatusNotFound
		}
		if errors.Is(err, usecaseErrors.ErrQuestionAlreadyAnswered) {
			httpStatus = http.StatusNotAcceptable
		}
		if errors.Is(err, usecaseErrors.ErrQuestionSessionNotRunning) {
			httpStatus = http.StatusNotAcceptable
		}
		if errors.Is(err, usecaseErrors.ErrQuestionNotOwned) {
			httpStatus = http.StatusForbidden
		}
		if errors.Is(err, usecaseErrors.ErrUnexpected) {
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

	err := usecases.Upvote(questionId, *userContext)

	httpStatus := http.StatusOK
	if err != nil {
		if errors.Is(err, usecaseErrors.ErrQuestionNotFound) {
			httpStatus = http.StatusNotFound
		}
		if errors.Is(err, usecaseErrors.ErrQuestionAlreadyAnswered) {
			httpStatus = http.StatusNotAcceptable
		}
		if errors.Is(err, usecaseErrors.ErrUserAlreadyVoted) {
			httpStatus = http.StatusNotAcceptable
		}
		if errors.Is(err, usecaseErrors.ErrQuestionSessionNotRunning) {
			httpStatus = http.StatusNotAcceptable
		}
		if errors.Is(err, usecaseErrors.ErrUnexpected) {
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

	err := usecases.UndoVote(questionId, *userContext)

	httpStatus := http.StatusOK
	if err != nil {
		if errors.Is(err, usecaseErrors.ErrQuestionNotFound) {
			httpStatus = http.StatusNotFound
		}
		if errors.Is(err, usecaseErrors.ErrQuestionAlreadyAnswered) {
			httpStatus = http.StatusNotAcceptable
		}
		if errors.Is(err, usecaseErrors.ErrUserHasNotVoted) {
			httpStatus = http.StatusNotAcceptable
		}
		if errors.Is(err, usecaseErrors.ErrQuestionSessionNotRunning) {
			httpStatus = http.StatusNotAcceptable
		}
		if errors.Is(err, usecaseErrors.ErrUnexpected) {
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

	err := usecases.Answer(questionId)

	httpStatus := http.StatusOK
	if err != nil {
		if errors.Is(err, usecaseErrors.ErrQuestionNotFound) {
			httpStatus = http.StatusNotFound
		}
		if errors.Is(err, usecaseErrors.ErrQuestionAlreadyAnswered) {
			httpStatus = http.StatusNotAcceptable
		}
		if errors.Is(err, usecaseErrors.ErrQuestionSessionNotRunning) {
			httpStatus = http.StatusNotAcceptable
		}
		if errors.Is(err, usecaseErrors.ErrUnexpected) {
			httpStatus = http.StatusBadRequest
		}
	}

	if err != nil {
		c.JSON(int(httpStatus), gin.H{
			"error": err.Error(),
		})
	}
}
