package services

import (
	"github.com/gin-gonic/gin"
)

type QuestionService interface {
	UndoVote(c *gin.Context)
	Answer(c *gin.Context)
	GetSession(c *gin.Context)
}
