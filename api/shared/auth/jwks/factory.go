package jwks

import (
	"voting/internal/env"
)

func Init() {
	if !env.Env.UseMockJwks {
		create()
	} else {
		Mock()
	}
}
