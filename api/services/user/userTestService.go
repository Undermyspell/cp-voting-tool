package services

import (
	"net/http"
	"sse/dtos"
	"sse/internal/mocks"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TestUserService struct{}

// GetContributor         godoc
// @Security 	 JWT
// @Summary      Gets a test user without role
// @Description  Gets a test user without role
// @Tags         User
// @Produce      json
// @Param        question  body dtos.NewTestUserDto  true  "User JSON"
// @Success      200
// @Failure      404
// @Router       /user/test/contributor [post]
func (service *TestUserService) GetContributor(c *gin.Context) {
	var userDto dtos.NewTestUserDto

	err := c.BindJSON(&userDto)

	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	token := mocks.GetUserToken(userDto.FirstName, userDto.LastName)

	c.JSON(http.StatusOK, struct{ Token string }{Token: token})
}

// GetAdmin         godoc
// @Security 	 JWT
// @Summary      Gets a test user with admin role
// @Description  Gets a test user with admin role
// @Tags         User
// @Produce      json
// @Param        question  body dtos.NewTestUserDto  true  "User JSON"
// @Success      200
// @Failure      404
// @Router       /user/test/admin [post]
func (service *TestUserService) GetAdmin(c *gin.Context) {
	var userDto dtos.NewTestUserDto

	err := c.BindJSON(&userDto)

	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	token := mocks.GetAdminUserToken(userDto.FirstName, userDto.LastName)

	c.JSON(http.StatusOK, struct{ Token string }{Token: token})
}

// GetSessionAdmin         godoc
// @Security 	 JWT
// @Summary      Gets a test user with session admin role
// @Description  Gets a test user with session admin role
// @Tags         User
// @Produce      json
// @Param        question  body dtos.NewTestUserDto  true  "User JSON"
// @Success      200
// @Failure      404
// @Router       /user/test/admin [post]
func (service *TestUserService) GetSessionAdmin(c *gin.Context) {
	var userDto dtos.NewTestUserDto

	err := c.BindJSON(&userDto)

	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusBadRequest, "cant parse request")
		return
	}

	token := mocks.GetSessionAdminUserToken(userDto.FirstName, userDto.LastName)

	c.JSON(http.StatusOK, struct{ Token string }{Token: token})
}
