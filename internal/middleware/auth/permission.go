package auth

import (
	"FastGo/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// GetMiddleware 根据权限等级获取对应的中间件
func GetMiddleware(level int) gin.HandlerFunc {
	switch level {
	case jwt.REQUIRE_AUTH:
		return TokenAuth()
	case jwt.TEAM_MEMBER:
		return TeamAuth(jwt.TEAM_MEMBER)
	case jwt.TEAM_ADMIN:
		return TeamAuth(jwt.TEAM_ADMIN)
	case jwt.TEAM_OWNER:
		return TeamAuth(jwt.TEAM_OWNER)
	default:
		return func(c *gin.Context) { c.Next() }
	}
}
