package voting_repositories

import (
	"context"
	"flag"
	"voting/internal/env"
	"voting/internal/helper"

	"github.com/nitishm/go-rejson/v4"
	"github.com/redis/go-redis/v9"
)

func NewInMemory() *InMemory {
	return &InMemory{}
}

func NewRedis() *Redis {
	redisEndpoint := env.Env.RedisEndpoint
	redisPassword := env.Env.RedisPassword
	generateRedisStorageRootKey := env.Env.GenerateRedisStorageRootKey

	redisRk := DefaultStorageRootKey

	if generateRedisStorageRootKey {
		redisRk = helper.GetRandomId(15)
	}

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
