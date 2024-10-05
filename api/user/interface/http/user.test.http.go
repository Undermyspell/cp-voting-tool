package user_http

import (
	"net/http"
	user_usecases "voting/user/usecases"
	user_dtos "voting/user/usecases/_dtos"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetContributor(c *gin.Context) {
	var userDto user_dtos.NewTestUserDto

	err := c.BindJSON(&userDto)

	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	token := user_usecases.GetContributorToken(userDto.FirstName, userDto.LastName)

	c.JSON(http.StatusOK, struct{ Token string }{Token: token})
}

func GetAdmin(c *gin.Context) {
	var userDto user_dtos.NewTestUserDto

	err := c.BindJSON(&userDto)

	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	token := user_usecases.GetAdminUserToken(userDto.FirstName, userDto.LastName)

	c.JSON(http.StatusOK, struct{ Token string }{Token: token})
}

func GetSessionAdmin(c *gin.Context) {
	var userDto user_dtos.NewTestUserDto

	err := c.BindJSON(&userDto)

	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	token := user_usecases.GetSessionAdminUserToken(userDto.FirstName, userDto.LastName)

	c.JSON(http.StatusOK, struct{ Token string }{Token: token})
}
