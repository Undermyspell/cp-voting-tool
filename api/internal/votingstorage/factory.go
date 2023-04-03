package votingstorage

import (
	"flag"

	goredisv8 "github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

func NewInMemory() *InMemory {
	return &InMemory{}
}

func NewRedis() *Redis {
	var addr = flag.String("Server", "redis-11436.c250.eu-central-1-1.ec2.cloud.redislabs.com:11436", "Redis server address")
	rh := rejson.NewReJSONHandler()

	cli := goredisv8.NewClient(&goredisv8.Options{
		Addr:     *addr,
		Password: "ZfSUj5EMRm7WilIOZnw0kORB8dWizA9A",
		DB:       0,
	})
	rh.SetGoRedisClient(cli)
	return &Redis{
		redisHandler: rh,
		redisClient:  cli,
	}
}
