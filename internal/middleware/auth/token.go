package auth

import (
	"FastGo/pkg/jwt"
	"FastGo/pkg/response"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

// 自定义错误
var (
	ErrTokenExpired     = errors.New(response.GetMessage(response.TokenExpired))
	ErrTokenNotValidYet = errors.New(response.GetMessage(response.TokenNotValidYet))
	ErrInvalidToken     = errors.New(response.GetMessage(response.InvalidToken))
)

// TokenAuth Token认证中间件
func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.NewResult(c).Fail(response.Unauthorized)
			c.Abort()
			return
		}

		// 移除Bearer前缀
		token = strings.TrimPrefix(token, "Bearer ")

		// 解析 token
		j := jwt.New("your-secret-key") // TODO: 从配置文件获取密钥
		claims, err := j.ParseToken(token)
		if err != nil {
			switch err {
			case jwt.ErrTokenExpired:
				response.NewResult(c).Fail(response.Unauthorized)
			case jwt.ErrTokenNotValidYet:
				response.NewResult(c).Fail(response.TokenNotValidYet)
			default:
				response.NewResult(c).Fail(response.InvalidToken)
			}
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
