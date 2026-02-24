package dao

import (
	"bluebell/logger"
	"bluebell/settings"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewRedis(lc fx.Lifecycle) *redis.Client {
	rcfg := settings.GlobalConfig.RedisConfig
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", rcfg.Host, rcfg.Port),
		Password:     rcfg.Password,
		DB:           rcfg.DB,
		PoolSize:     rcfg.PoolSize,
		MinIdleConns: rcfg.MinIdleConns,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logger.L(nil).Panic("failed to connect to redis", zap.Error(err))
		panic(err)
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			zap.L().Debug("closing redis client")
			err := client.Close()
			if err != nil {
				zap.L().Error("failed to close redis client", zap.Error(err))
				return err
			}
			zap.L().Debug("redis client closed")
			return nil
		},
	})
	return client
}
