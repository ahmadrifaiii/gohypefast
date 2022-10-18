package config

import "github.com/go-redis/redis/v8"

type Configuration struct {
	RedisConnect *redis.Client
}
