package env

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type EnvConfig struct {
	UseMockJwks bool
}

var Env EnvConfig

func Init() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	Env = EnvConfig{
		UseMockJwks: viper.GetBool("USE_MOCK_JWKS"),
	}

	logrus.Infof("Found env var: %s", Env)
}
