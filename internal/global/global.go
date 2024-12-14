package global

import (
	"FastGo/config"

	"sync"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config        *config.Options
	Engine        *gin.Engine
	dbInstance    *gorm.DB
	dbOnce        sync.Once
	redisInstance *redis.Client
	redisOnce     sync.Once
	Log           *zap.Logger
)

// GetDB 返回数据库实例
func GetDB() *gorm.DB {
	return dbInstance
}

// SetDB 设置数据库实例
func SetDB(db *gorm.DB) {
	dbOnce.Do(func() {
		dbInstance = db
	})
}

// GetRedis 返回 Redis 客户端实例
func GetRedis() *redis.Client {
	return redisInstance
}

// SetRedis 设置 Redis 客户端实例
func SetRedis(client *redis.Client) {
	redisOnce.Do(func() {
		redisInstance = client
	})
}
