package core

import (
	"context"

	"backend-go/pkg/config"
	"backend-go/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RDB *redis.Client

func InitRedis() *redis.Client {
	cfg := config.C.Redis

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Log.Fatal("failed to connect redis", zap.Error(err))
	}

	RDB = rdb
	logger.Info("Redis connected successfully")
	return rdb
}
