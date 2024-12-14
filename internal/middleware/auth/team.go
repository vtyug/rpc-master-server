package auth

import (
	"FastGo/pkg/response"

	"github.com/gin-gonic/gin"
)

// TeamAuth 团队权限验证中间件
func TeamAuth(requiredLevel uint8) gin.HandlerFunc {
	return func(c *gin.Context) {
		TokenAuth()(c)
		if c.IsAborted() {
			return
		}

		role, exists := c.Get("role")
		if !exists {
			response.NewResult(c).Fail(response.Unauthorized)
			c.Abort()
			return
		}

		userRole := role.(uint8)
		if userRole < requiredLevel {
			response.NewResult(c).Fail(response.Forbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
