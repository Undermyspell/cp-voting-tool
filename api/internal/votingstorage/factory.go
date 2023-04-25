package votingstorage

import (
	"flag"
	"sse/internal/env"

	goredisv8 "github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

func NewInMemory() *InMemory {
	return &InMemory{}
}

func NewRedis() *Redis {
	redisEndpoint := env.Env.RedisEndpoint
	redisPassword := env.Env.RedisPassword

	var addr = flag.String("Server", redisEndpoint, "Redis server address")
	rh := rejson.NewReJSONHandler()

	cli := goredisv8.NewClient(&goredisv8.Options{
		Addr:     *addr,
		Password: redisPassword,
		DB:       0,
	})
	rh.SetGoRedisClient(cli)
	return &Redis{
		redisHandler: rh,
		redisClient:  cli,
	}
}