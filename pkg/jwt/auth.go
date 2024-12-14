package jwt

// 权限等级
const (
	PUBLIC       = 1 // 公开接口
	REQUIRE_AUTH = 2 // 需要认证
	TEAM_MEMBER  = 3 // 团队成员
	TEAM_ADMIN   = 4 // 团队管理员
	TEAM_OWNER   = 5 // 团队拥有者
)

// GetAuthName 获取权限名称
func GetAuthName(level int) string {
	switch level {
	case PUBLIC:
		return "公开"
	case REQUIRE_AUTH:
		return "认证"
	case TEAM_MEMBER:
		return "团队成员"
	case TEAM_ADMIN:
		return "团队管理员"
	case TEAM_OWNER:
		return "团队拥有者"
	default:
		return "未知"
	}
}
