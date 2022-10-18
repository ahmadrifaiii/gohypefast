package redisconf

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"hypefast.io/services/config/env"
	"hypefast.io/services/pkg/utils/logging"
)

var (
	master     *redis.Client
	lockMaster sync.Mutex
	ctx        = context.Background()
)

func GetMasterRedis() *redis.Client {
	lockMaster.Lock()
	defer lockMaster.Unlock()

	if master == nil {
		master = newConnection()
	}

	logging.Info("info", zap.String("message", "redis connection is connect"))
	return master
}

func newConnection() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", env.Conf.RedisHost, env.Conf.RedisPort),
		Password:     env.Conf.RedisPassword,
		PoolTimeout:  20 * time.Second,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis can't connect : %v\n", err)
	}

	return client
}

func RedisConnect() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", env.Conf.RedisHost, env.Conf.RedisPort),
		Password:     env.Conf.RedisPassword,
		PoolTimeout:  20 * time.Second,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis can't connect : %v\n", err)
	}

	return client
}
