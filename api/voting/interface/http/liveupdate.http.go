package voting_http

import (
	"net/http"
	"voting/internal/env"
	"voting/shared/auth/middleware"
	voting_ws "voting/voting/interface/ws"

	"github.com/centrifugal/centrifuge"
	"github.com/gin-gonic/gin"
)

func CentrifugoHandler() gin.HandlerFunc {
	return gin.WrapH(middleware.CentrifugeAnonymousAuth(centrifuge.NewWebsocketHandler(voting_ws.GetHandler(), centrifuge.WebsocketConfig{
		CheckOrigin: func(r *http.Request) bool {
			originHeader := r.Header.Get("Origin")
			if originHeader == "" {
				return true
			}
			for _, allowedOrigin := range env.Env.AllowedOrigins {
				if originHeader == allowedOrigin {
					return true
				}
			}
			return false
		},
	})))
}
