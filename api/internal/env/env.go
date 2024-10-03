package env

import (
	"strings"

	"github.com/spf13/viper"
)

type Storage string

const (
	Redis    Storage = "redis"
	Postgres Storage = "postgres"
	InMemory Storage = "inmemory"
)

const (
	USE_MOCK_JWKS                string = "USE_MOCK_JWKS"
	JWKS_URL                     string = "JWKS_URL"
	STORAGE                      string = "STORAGE"
	REDIS_ENPOINT                string = "REDIS_ENDPOINT_SECRET"
	REDIS_PASSWORD               string = "REDIS_PASSWORD_SECRET"
	POSTGRESQL_CONNECTION_STRING string = "POSTGRESQL_CONNECTION_STRING_SECRET"
	ALLOWED_ORIGINS              string = "ALLOWED_ORIGINS"
	SESSION_SECRET               string = "SESSION_SECRET"
	AUTH_AZURE_CLIENT_ID         string = "AUTH_AZURE_CLIENT_ID"
	AUTH_AZURE_CLIENT_SECRET     string = "AUTH_AZURE_CLIENT_SECRET"
	AUTH_AZURE_TENANT_ID         string = "AUTH_AZURE_TENANT_ID"
	AUTH_REDIRECT_URL            string = "AUTH_REDIRECT_URL"
)

type EnvConfig struct {
	UseMockJwks                bool
	Storage                    Storage
	JwksUrl                    string
	RedisEndpoint              string
	RedisPassword              string
	AllowedOrigins             []string
	PostgresqlConnectionString string
	SessionSecret              string
	AuthAzureClientId          string
	AuthAzureClientSecret      string
	AuthAzureTenantId          string
	AuthRedirectUrl            string
}

var Env EnvConfig

func Init() {
	viper.AutomaticEnv()
	viper.SetConfigFile("env.yaml")
	viper.ReadInConfig()

	Env = EnvConfig{
		UseMockJwks:                viper.GetBool(USE_MOCK_JWKS),
		JwksUrl:                    viper.GetString(JWKS_URL),
		Storage:                    Storage(viper.GetString(STORAGE)),
		RedisEndpoint:              viper.GetString(REDIS_ENPOINT),
		RedisPassword:              viper.GetString(REDIS_PASSWORD),
		AllowedOrigins:             strings.Split(viper.GetString(ALLOWED_ORIGINS), ","),
		PostgresqlConnectionString: viper.GetString(POSTGRESQL_CONNECTION_STRING),
		SessionSecret:              viper.GetString(SESSION_SECRET),
		AuthAzureClientId:          viper.GetString(AUTH_AZURE_CLIENT_ID),
		AuthAzureClientSecret:      viper.GetString(AUTH_AZURE_CLIENT_SECRET),
		AuthAzureTenantId:          viper.GetString(AUTH_AZURE_TENANT_ID),
		AuthRedirectUrl:            viper.GetString(AUTH_REDIRECT_URL),
	}
}
