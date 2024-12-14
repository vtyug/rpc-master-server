package auth

import (
	"FastGo/pkg/response"
	"fmt"

	"github.com/gin-gonic/gin"
)

// ApiKeyAuth API Key认证中间件
func ApiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			response.NewResult(c).Fail(response.Unauthorized)
			c.Abort()
			return
		}

		// 验证API Key
		userID, err := validateApiKey(apiKey)
		if err != nil {
			response.NewResult(c).FailWithError(response.Unauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

// validateApiKey 验证API Key
func validateApiKey(apiKey string) (uint, error) {
	// TODO: 实现API Key验证
	fmt.Println("apiKey", apiKey)
	return 0, nil
}
