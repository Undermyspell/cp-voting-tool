package bff

import (
	"fmt"
	"strconv"
	"time"
	"voting/bff/templates/components"
	"voting/bff/templates/pages"
	"voting/shared/helper/httputils"
	voting_usecases "voting/voting/usecases"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewQuestion(c *gin.Context) {
	component := components.NewQuestionModal()
	component.Render(c.Request.Context(), c.Writer)
}

func SaveQuestion(c *gin.Context) {
	question := c.PostForm("question")
	anonymous, _ := strconv.ParseBool(c.DefaultPostForm("anonymous", "false"))

	dto := voting_usecases.NewQuestionDto{
		Text:      question,
		Anonymous: anonymous,
	}

	logrus.Info(fmt.Sprintf("%s:%v", question, anonymous))

	sessions := sessions.Default(c)
	token := sessions.Get("token").(string)
	httputils.Post("http://:3333/api/v1/question/new", map[string]string{
		"Authorization": "Bearer " + token,
	}, dto)

	time.Sleep(2 * time.Second)
	component := components.SuccessToast()
	component.Render(c.Request.Context(), c.Writer)
}

func Home(c *gin.Context) {
	sessions := sessions.Default(c)
	token := sessions.Get("token").(string)

	questions, _ := httputils.Get[[]voting_usecases.QuestionDto]("http://:3333/api/v1/question/session", map[string]string{
		"Authorization": "Bearer " + token,
	})

	component := pages.Main("Home Page :)", "Welcome to the Home Page :)!", *questions)
	component.Render(c.Request.Context(), c.Writer)
}
