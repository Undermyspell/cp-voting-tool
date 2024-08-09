package bff

import (
	"voting/bff/templates/components"
	"voting/bff/templates/pages"
	"voting/shared/helper/httputils"
	voting_usecases "voting/voting/usecases"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func NewQuestion(c *gin.Context) {
	component := components.NewQuestionModal()
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
