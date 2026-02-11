package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

var Ctx = context.Background()

func ConnectRedis(cfg *Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Redis connection error: ", err)
	}
	log.Println("Successfully connected to redis")
	return rdb
}
