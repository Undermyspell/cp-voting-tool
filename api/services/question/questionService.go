package services

import (
	"github.com/gin-gonic/gin"
)

type QuestionService interface {
	Add(c *gin.Context)
	Update(c *gin.Context)
	Upvote(c *gin.Context)
	UndoVote(c *gin.Context)
	Delete(c *gin.Context)
	Answer(c *gin.Context)
	Stop(c *gin.Context)
	Start(c *gin.Context)
	GetSession(c *gin.Context)
}
