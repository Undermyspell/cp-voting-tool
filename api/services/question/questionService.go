package services

import (
	"github.com/gin-gonic/gin"
)

type QuestionService interface {
	Answer(c *gin.Context)
	GetSession(c *gin.Context)
}
