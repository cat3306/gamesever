package thirdmodule

import (
	"github.com/go-redis/redis"
)

var (
	Cache *redis.Client
)

func InitCache() error {
	Cache = redis.NewClient(&redis.Options{
		Addr:     "",
		Password: "",
		DB:       0,
	})

	return nil
}
