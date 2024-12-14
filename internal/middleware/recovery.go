package middleware

import (
	"FastGo/internal/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.Log.Error("服务器内部错误",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
				)

				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "服务器内部错误",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
