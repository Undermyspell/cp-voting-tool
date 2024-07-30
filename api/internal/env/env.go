package env

import (
	"github.com/spf13/viper"
)

const (
	USE_MOCK_JWKS                string = "USE_MOCK_JWKS"
	JWKS_URL                     string = "JWKS_URL"
	VOTING_STORAGE_IN_MEMORY     string = "VOTING_STORAGE_IN_MEMORY"
	REDIS_ENPOINT                string = "REDIS_ENDPOINT_SECRET"
	REDIS_PASSWORD               string = "REDIS_PASSWORD_SECRET"
	POSTGRESQL_CONNECTION_STRING string = "POSTGRESQL_CONNECTION_STRING_SECRET"
	ALLOWED_ORIGIN               string = "ALLOWED_ORIGIN"
)

type EnvConfig struct {
	UseMockJwks                bool
	VotingStorageInMemory      bool
	JwksUrl                    string
	RedisEndpoint              string
	RedisPassword              string
	AllowedOrigin              string
	PostgresqlConnectionString string
}

var Env EnvConfig

func Init() {
	viper.AutomaticEnv()
	viper.SetConfigFile("env.yaml")
	viper.ReadInConfig()

	Env = EnvConfig{
		UseMockJwks:                viper.GetBool(USE_MOCK_JWKS),
		JwksUrl:                    viper.GetString(JWKS_URL),
		VotingStorageInMemory:      viper.GetBool(VOTING_STORAGE_IN_MEMORY),
		RedisEndpoint:              viper.GetString(REDIS_ENPOINT),
		RedisPassword:              viper.GetString(REDIS_PASSWORD),
		AllowedOrigin:              viper.GetString(ALLOWED_ORIGIN),
		PostgresqlConnectionString: viper.GetString(POSTGRESQL_CONNECTION_STRING),
	}
}
