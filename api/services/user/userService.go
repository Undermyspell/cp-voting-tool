package services

import (
	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetContributor(c *gin.Context)
	GetAdmin(c *gin.Context)
	GetSessionAdmin(c *gin.Context)
}
