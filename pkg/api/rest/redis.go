package rest

import (
	"github.com/go-redis/redis/v8"
)

// var rdb *redis.Client

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}
