package votingstorage

import (
	"flag"

	goredis "github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

func NewInMemory() *InMemory {
	return &InMemory{}
}

func NewRedis() *Redis {
	var addr = flag.String("Server", "redis-11436.c250.eu-central-1-1.ec2.cloud.redislabs.com:11436", "Redis server address")
	rh := rejson.NewReJSONHandler()
	cli := goredis.NewClient(&goredis.Options{
		Addr:     *addr,
		Password: "ZfSUj5EMRm7WilIOZnw0kORB8dWizA9A",
		DB:       0,
	})
	rh.SetGoRedisClient(cli)
	return &Redis{
		redisHandler: rh,
		goRedisCli:   cli,
	}
}
