package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"FastGo/internal/global"
	"FastGo/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Server struct {
	app    *App
	engine *gin.Engine
	srv    *http.Server
}

func NewServer(app *App) *Server {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	// 使用中间件
	engine.Use(
		middleware.Logger(),   // 日志中间件
		middleware.Recovery(), // 恢复中间件
		middleware.Cors(),     // 跨域中间件
		// jwt.JWT(),       // JWT 中间件（按需启用）
	)

	global.Engine = engine

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", 8080),
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &Server{
		app:    app,
		engine: engine,
		srv:    srv,
	}
}

// InitRoutes 初始化路由
func (s *Server) InitRoutes() {
	// 健康检查
	s.engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

// Start 启动服务器
func (s *Server) Start() error {
	global.Log.Info("服务器启动中...",
		zap.String("addr", s.srv.Addr),
	)

	// 打印所有注册的路由
	routes := s.engine.Routes()
	if len(routes) > 0 {
		global.Log.Info("已注册的路由列表:")
		for _, route := range routes {
			global.Log.Info("→",
				zap.String("method", route.Method),
				zap.String("path", route.Path),
			)
		}
	} else {
		global.Log.Warn("没有找到任何注册的路由!")
	}

	// 初始化路由
	s.InitRoutes()

	// 保存到 app 中，方便优雅关闭
	s.app.Server = s.srv

	return s.srv.ListenAndServe()
}

// Shutdown 关闭服务器
func (s *Server) Shutdown(ctx context.Context) error {
	global.Log.Info("服务器正在关闭...")
	return s.srv.Shutdown(ctx)
}
