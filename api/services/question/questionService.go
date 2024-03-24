package services

import (
	"github.com/gin-gonic/gin"
)

type QuestionService interface {
	GetSession(c *gin.Context)
}
