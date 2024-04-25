package voting_repositories

import (
	"context"
	"flag"
	"voting/internal/env"

	"github.com/nitishm/go-rejson/v4"
	"github.com/redis/go-redis/v9"
)

func NewInMemory() *InMemory {
	return &InMemory{}
}

func NewRedis() *Redis {
	redisEndpoint := env.Env.RedisEndpoint
	redisPassword := env.Env.RedisPassword

	redisRk := DefaultStorageRootKey

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
