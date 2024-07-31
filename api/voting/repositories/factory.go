package voting_repositories

import (
	"context"
	"flag"
	"voting/internal/env"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nitishm/go-rejson/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func NewInMemory() *InMemory {
	return &InMemory{}
}

func NewRedis() *Redis {
	redisEndpoint := env.Env.RedisEndpoint
	redisPassword := env.Env.RedisPassword

	redisRk := DefaultVotingSessionRootKey

	var addr = flag.String("Server", redisEndpoint, "Redis server address")
	rh := rejson.NewReJSONHandler()

	cli := redis.NewClient(&redis.Options{
		Addr:     *addr,
		Password: redisPassword,
		DB:       0,
	})
	rh.SetGoRedisClientWithContext(context.Background(), cli)

	return &Redis{
		redisHandler: rh,
		redisClient:  cli,
		redisRootKey: redisRk,
	}
}

func NewPostgresql() *Postgresql {
	postgresqlConnectionString := env.Env.PostgresqlConnectionString

	poolConfig, err := pgxpool.ParseConfig(postgresqlConnectionString)

	if err != nil {
		logrus.Fatalf("Unable to create db pool config %v\n", err)
	}

	poolConfig.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)

	if err != nil {
		logrus.Fatalf("Unable to connect to PostgreSQL database: %v\n", err)
	}

	sessionKey := DefaultVotingSessionRootKey

	return &Postgresql{
		pool:      pool,
		sessionId: sessionKey,
	}
}
