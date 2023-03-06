package env

import (
	"github.com/spf13/viper"
)

const (
	USE_MOCK_JWKS string = "USE_MOCK_JWKS"
	JWKS_URL      string = "JWKS_URL"
)

type EnvConfig struct {
	UseMockJwks bool
	JwksUrl     string
}

var Env EnvConfig

func Init() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	Env = EnvConfig{
		UseMockJwks: viper.GetBool(USE_MOCK_JWKS),
		JwksUrl:     viper.GetString(JWKS_URL),
	}
}
