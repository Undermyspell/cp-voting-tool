package middleware

import (
	"log"
	"net/http"
	"sse/internal/jwks"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")

		splitted := strings.Split(bearerToken, " ")

		if len(splitted) != 2 && splitted[0] != "Bearer" {
			log.Println("Not a valid bearer token")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		jwtB64 := splitted[1]
		jwks := jwks.Jwks

		token, err := jwt.Parse(jwtB64, jwks.Keyfunc)
		if err != nil {
			log.Printf("Failed to parse the JWT.\nError: %s", err.Error())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Printf("The token is not valid.")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		log.Println("The token is valid.")

		c.Next()
	}
}
