package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/nishant1479/Microservice-Backend/config"
	"github.com/redis/go-redis"
)

func NewRedisClient(cfg *config.Config) *redis.Client{
	redisHost := cfg.Redis.RedisAddr

	if redisHost == ""{
		redisHost = ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: 			redisHost,
		MinIdleCons:	cfg.Redis.MinIdleConn,
		PoolSize:		cfg.Redis.PoolSize,
		PoolTimeout:	time.Duration(cfg.Redis.PoolTimeout) * time.Second,
		Password:		cfg.Redis.Password,
		DB:				cfg.Redis.DB,
	})

	return client
}