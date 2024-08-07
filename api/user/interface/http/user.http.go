package user_http

import (
	"net/http"
	shared_models "voting/shared/models"
	user_dtos "voting/user/usecases/_dtos"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetAuthenticatedUser(c *gin.Context) {
	session := sessions.Default(c)

	token, _ := session.Get("token").(string)

	userContext, _ := shared_models.GetUserContextFromToken(token)

	userDto := user_dtos.UserDto{
		Email: userContext.Email,
		Name:  userContext.Name,
		Token: token,
		Role:  userContext.Role,
	}

	c.JSON(http.StatusOK, userDto)
}
