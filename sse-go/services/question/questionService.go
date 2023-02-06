package services

import "github.com/gin-gonic/gin"

type QuestionService interface {
	AddQuestion(c *gin.Context)
	UpvoteQuestion(c *gin.Context)
	Answer(c *gin.Context)
	Reset(c *gin.Context)
}
