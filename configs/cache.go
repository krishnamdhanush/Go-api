package configs

import (
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func CreateRedisClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     EnvRedis(),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	RedisClient = client
}
