package config

import (
	"github.com/go-redis/redis"
)

type RedisServer redis.Options

func NewRedis() (rdb *redis.Client){
	return redis.NewClient((*redis.Options)(&RedisServer{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}))
}






