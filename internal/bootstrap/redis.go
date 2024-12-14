package bootstrap

import (
	"FastGo/internal/global"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (app *App) setupRedis() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", app.Config.Redis.Host, app.Config.Redis.Port),
		Password: app.Config.Redis.Password,
		DB:       app.Config.Redis.DB,
		PoolSize: app.Config.Redis.PoolSize,
	})

	// 测试连接
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		global.Log.Error("连接Redis失败", zap.Error(err))
		return fmt.Errorf("连接Redis失败: %w", err)
	}

	// 使用单例模式设置 Redis 客户端实例
	global.SetRedis(rdb)

	global.Log.Info("Redis 连接成功")

	return nil
}
