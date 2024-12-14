package handler

import (
	"FastGo/internal/global"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CommonHandler 包含全局对象的通用结构体
type CommonHandler struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Logger *zap.Logger
}

// NewCommonHandler 创建并返回一个 CommonHandler 实例
func NewCommonHandler() *CommonHandler {
	return &CommonHandler{
		DB:     global.GetDB(),
		Redis:  global.GetRedis(),
		Logger: global.Log,
	}
}
