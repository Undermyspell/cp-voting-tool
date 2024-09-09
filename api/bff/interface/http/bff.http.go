package bff

import (
	"fmt"
	"net/http"
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

func UpdateQuestion(c *gin.Context) {
	questionId := c.Param("id")

	sessions := sessions.Default(c)
	token := sessions.Get("token").(string)

	questions, _ := httputils.Get[[]voting_usecases.QuestionDto]("http://:3333/api/v1/question/session", map[string]string{
		"Authorization": "Bearer " + token,
	})

	var questionToEdit voting_usecases.QuestionDto
	for _, question := range *questions {
		if question.Id == questionId {
			questionToEdit = question
			break
		}
	}

	component := components.UpdateQuestionModal(questionToEdit)
	component.Render(c.Request.Context(), c.Writer)
}

func SaveNewQuestion(c *gin.Context) {
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
	component := components.SuccessToast("Frage erfolgreich erstellt")
	component.Render(c.Request.Context(), c.Writer)
}

func SaveUpdatedQuestion(c *gin.Context) {

	// Parse the form data (needed for PUT requests)
	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}
	id := c.PostForm("id")
	question := c.PostForm("question")
	anonymous, _ := strconv.ParseBool(c.DefaultPostForm("anonymous", "false"))

	dto := voting_usecases.UpdateQuestionDto{
		Id:        id,
		Text:      question,
		Anonymous: anonymous,
	}

	logrus.Info(fmt.Sprintf("%s:%v", question, anonymous))

	sessions := sessions.Default(c)
	token := sessions.Get("token").(string)
	httputils.Put("http://:3333/api/v1/question/update", map[string]string{
		"Authorization": "Bearer " + token,
	}, dto)

	time.Sleep(2 * time.Second)
	component := components.SuccessToast("Frage erfolgreich aktualisiert")
	component.Render(c.Request.Context(), c.Writer)
}

func DeleteQuestion(c *gin.Context) {
	questionId := c.Param("id")

	sessions := sessions.Default(c)
	token := sessions.Get("token").(string)
	httputils.Delete("http://:3333/api/v1/question/delete/"+questionId, map[string]string{
		"Authorization": "Bearer " + token,
	})

	time.Sleep(2 * time.Second)
	component := components.SuccessToast("Frage erfolgreich gelöscht")
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
