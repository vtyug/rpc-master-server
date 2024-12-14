package router

import (
	"FastGo/internal/global"
	"FastGo/internal/middleware/auth"
	"log"

	"github.com/gin-gonic/gin"
)

// RouteConfig 定义路由配置结构体
type RouteConfig struct {
	Method      string
	Path        string
	HandlerFunc gin.HandlerFunc
	AuthLevel   int
	Description string
	Group       string
}

// RouteRegistry 路由注册器
type RouteRegistry struct {
	routes []RouteConfig
}

// NewRouteRegistry 创建并返回一个新的路由注册器
func NewRouteRegistry() *RouteRegistry {
	return &RouteRegistry{
		routes: []RouteConfig{},
	}
}

// Register 注册一个新的路由
func (r *RouteRegistry) Register(method, group, path string, handlerFunc gin.HandlerFunc, authLevel int, description string) {
	if group == "" {
		panic("group is required")
	}
	fullPath := "/" + group + path
	r.routes = append(r.routes, RouteConfig{
		Method:      method,
		Path:        fullPath,
		HandlerFunc: handlerFunc,
		AuthLevel:   authLevel,
		Description: description,
		Group:       group,
	})
}

// SetupRoutes 注册所有路由
func (r *RouteRegistry) SetupRoutes(engine *gin.Engine) {
	for _, route := range r.routes {
		RegisterRoute(route.Method, route.Path, route.HandlerFunc, route.AuthLevel, route.Description)
	}
}

// RegisterRoute 注册路由并处理权限
func RegisterRoute(method string, path string, handler gin.HandlerFunc, authLevel int, description string) {
	authMiddleware := auth.GetMiddleware(authLevel)

	methodMap := map[string]func(string, ...gin.HandlerFunc) gin.IRoutes{
		"POST":   global.Engine.POST,
		"GET":    global.Engine.GET,
		"PUT":    global.Engine.PUT,
		"DELETE": global.Engine.DELETE,
	}

	// 检查方法是否存在于映射中
	if ginMethod, exists := methodMap[method]; exists {
		ginMethod(path, authMiddleware, handler)
	} else {
		log.Fatalf("不支持的方法: %s", method)
	}
}
