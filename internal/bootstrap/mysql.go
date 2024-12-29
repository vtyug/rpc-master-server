package bootstrap

import (
	"FastGo/internal/global"
	"FastGo/internal/model"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func (app *App) setupMySQL() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		app.Config.MySQL.Username,
		app.Config.MySQL.Password,
		app.Config.MySQL.Host,
		app.Config.MySQL.Port,
		app.Config.MySQL.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		global.Log.Error("MySQL 连接失败",
			zap.Error(err),
		)
		return fmt.Errorf("连接MySQL失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		global.Log.Error("获取 *sql.DB 失败",
			zap.Error(err),
		)
		return err
	}

	// 设置连接池
	sqlDB.SetMaxOpenConns(app.Config.MySQL.MaxOpenConns)
	sqlDB.SetMaxIdleConns(app.Config.MySQL.MaxIdleConns)

	// 使用单例模式设置数据库实例
	global.SetDB(db)

	global.Log.Info("MySQL 连接成功")

	return nil
}

// Migrate 执行数据库迁移
func Migrate() {
	if global.GetDB() == nil {
		global.Log.Error("数据库未初始化")
		return
	}

	err := global.GetDB().AutoMigrate(
		&model.User{},
		// &model.Team{},
		// &model.TeamMember{},
		// &model.TeamInvite{},
		// &model.Request{},
		// &model.Workspace{},
		&model.Collections{},
		// &model.Folder{},
		// &model.FolderClosure{},
	)
	if err != nil {
		global.Log.Error("数据库迁移失败", zap.Error(err))
	} else {
		global.Log.Info("数据库迁移成功")
	}
}
