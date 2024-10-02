package bff

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"voting/bff/templates/components"
	"voting/bff/templates/layouts"
	"voting/bff/templates/pages"
	"voting/shared/helper/httputils"
	shared_models "voting/shared/models"
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

	time.Sleep(2 * time.Second)
	sessions := sessions.Default(c)
	token := sessions.Get("token").(string)
	httputils.Post("http://:3333/api/v1/question/new", map[string]string{
		"Authorization": "Bearer " + token,
	}, dto)

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

	time.Sleep(2 * time.Second)

	httputils.Put("http://:3333/api/v1/question/update", map[string]string{
		"Authorization": "Bearer " + token,
	}, dto)

	component := components.SuccessToast("Frage erfolgreich aktualisiert")
	component.Render(c.Request.Context(), c.Writer)
}

func DeleteQuestion(c *gin.Context) {
	time.Sleep(time.Second * 2)
	questionId := c.Param("id")

	sessions := sessions.Default(c)
	token := sessions.Get("token").(string)
	httputils.Delete("http://:3333/api/v1/question/delete/"+questionId, map[string]string{
		"Authorization": "Bearer " + token,
	})

	component := components.SuccessToast("Frage erfolgreich gelöscht")
	component.Render(c.Request.Context(), c.Writer)
}

func UpvoteQuestion(c *gin.Context) {
	questionId := c.Param("id")

	sessions := sessions.Default(c)
	token := sessions.Get("token").(string)
	dto := "{}"
	httputils.Put("http://:3333/api/v1/question/upvote/"+questionId, map[string]string{
		"Authorization": "Bearer " + token,
	}, dto)
}

func UndoVoteQuestion(c *gin.Context) {
	questionId := c.Param("id")

	sessions := sessions.Default(c)
	token := sessions.Get("token").(string)
	dto := "{}"
	httputils.Put("http://:3333/api/v1/question/undovote/"+questionId, map[string]string{
		"Authorization": "Bearer " + token,
	}, dto)
}

func AnswerQuestion(c *gin.Context) {
	questionId := c.Param("id")

	sessions := sessions.Default(c)
	token := sessions.Get("token").(string)
	dto := "{}"
	httputils.Put("http://:3333/api/v1/question/answer/"+questionId, map[string]string{
		"Authorization": "Bearer " + token,
	}, dto)
}

func StartSession(c *gin.Context) {
	sessions := sessions.Default(c)
	token := sessions.Get("token").(string)

	dto := "{}"
	httputils.Post("http://:3333/api/v1/question/session/start", map[string]string{
		"Authorization": "Bearer " + token,
	}, dto)
}

func StopSession(c *gin.Context) {
	sessions := sessions.Default(c)
	token := sessions.Get("token").(string)
	dto := "{}"
	httputils.Post("http://:3333/api/v1/question/session/stop", map[string]string{
		"Authorization": "Bearer " + token,
	}, dto)
}

func QuestionSessionPage(c *gin.Context) {
	sessions := sessions.Default(c)
	token := sessions.Get("token").(string)
	activeSession, _ := strconv.ParseBool(c.Param("activeSession"))
	onlyUnanswered, _ := strconv.ParseBool(c.Param("onlyUnanswered"))

	userContext, err := shared_models.GetUserContextFromToken(token)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	questions := []voting_usecases.QuestionDto{}
	component := pages.QuestionsSession(questions, activeSession, onlyUnanswered, *userContext)
	component.Render(c.Request.Context(), c.Writer)
}

func DownloadSessionAsCsv(c *gin.Context) {
	sessions := sessions.Default(c)
	token := sessions.Get("token").(string)

	questions, statusCode := httputils.Get[[]voting_usecases.QuestionDto]("http://:3333/api/v1/question/session", map[string]string{
		"Authorization": "Bearer " + token,
	})

	if statusCode > 299 {
		c.Status(statusCode)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=questions.csv")
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(c.Writer)

	headers := []string{"Text", "Votes"}
	if err := writer.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV header"})
		return
	}

	for _, question := range *questions {
		row := []string{question.Text, strconv.Itoa(question.Votes)}
		if err := writer.Write(row); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV file rows"})
			return
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to flush CSV data"})
		return
	}
}

func Main(c *gin.Context) {
	sessions := sessions.Default(c)
	token := sessions.Get("token").(string)

	questions, statusCode := httputils.Get[[]voting_usecases.QuestionDto]("http://:3333/api/v1/question/session", map[string]string{
		"Authorization": "Bearer " + token,
	})

	activeSession := true
	if statusCode > 299 {
		activeSession = false
	}

	if questions == nil {
		questions = &[]voting_usecases.QuestionDto{}
	}

	userContext, err := shared_models.GetUserContextFromToken(token)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	component := pages.QuestionsSession(*questions, activeSession, true, *userContext)

	layouts.Default(component, "cp voting tool").Render(c.Request.Context(), c.Writer)
}
