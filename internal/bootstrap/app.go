package bootstrap

import (
	"FastGo/config"
	"FastGo/pkg/validator"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"FastGo/internal/global"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// App 应用结构体
type App struct {
	Config *config.Options
	DB     *gorm.DB
	Redis  *redis.Client
	Server *http.Server
	srv    *Server
}

// Setup 初始化应用
func Setup() (*App, error) {

	// 1. 初始化配置
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("load config failed: %v", err)
	}
	global.Config = cfg

	app := &App{
		Config: cfg,
	}

	// 2. 初始化日志
	if err := setupLogger(); err != nil {
		return nil, fmt.Errorf("setup logger failed: %v", err)
	}

	// 3. 初始化 MySQL
	if err := app.setupMySQL(); err != nil {
		return nil, fmt.Errorf("setup mysql failed: %v", err)
	}

	// 4. 初始化 Redis
	if err := app.setupRedis(); err != nil {
		return nil, fmt.Errorf("setup redis failed: %v", err)
	}

	// 5. 初始化 Server
	server := NewServer(app)
	app.Server = server.srv
	app.srv = server

	// 6. 初始化验证器
	validator.Setup()

	// 7. 设置路由
	SetupRoutes()

	return app, nil
}

// Run 运行应用
func (app *App) Run() error {
	// 启动服务器
	go func() {
		if err := app.srv.Start(); err != nil {
			if err != http.ErrServerClosed {
				global.Log.Error("服务器启动失败", zap.Error(err))
			}
		}
	}()

	// 等待信号
	app.GracefulShutdown()
	return nil
}

// GracefulShutdown 优雅退出
func (app *App) GracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	global.Log.Info("开始优雅关闭...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if app.Server != nil {
		global.Log.Info("正在关闭 HTTP 服务器...")
		if err := app.Server.Shutdown(ctx); err != nil {
			global.Log.Error("HTTP 服务器关闭错误",
				zap.Error(err),
			)
		}
	}

	if app.DB != nil {
		global.Log.Info("正在关闭数据库连接...")
		sqlDB, err := app.DB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				global.Log.Error("数据库关闭错误",
					zap.Error(err),
				)
			}
		}
	}

	if app.Redis != nil {
		global.Log.Info("正在关闭 Redis 连接...")
		if err := app.Redis.Close(); err != nil {
			global.Log.Error("Redis 关闭错误",
				zap.Error(err),
			)
		}
	}

	global.Log.Info("服务器已成功关闭")
}
