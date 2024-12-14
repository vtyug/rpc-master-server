package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Options struct {
	MySQL MySQLOptions `json:"mysql"`
	Redis RedisOptions `json:"redis"`
	Log   LogConfig    `json:"log"`
	Jwt   JWT          `json:"jwt"`
}

func New() (*Options, error) {
	// 获取当前文件的路径
	_, b, _, _ := runtime.Caller(0)
	// 获取项目根目录（config.go 在 config 目录下）
	projectRoot := filepath.Dir(filepath.Dir(b))

	viper.SetConfigName("config") // 配置文件名称
	viper.SetConfigType("yaml")   // 配置文件类型

	// 添加配置文件搜索路径
	viper.AddConfigPath(filepath.Join(projectRoot, "config")) // 项目根目录下的 config 目录
	viper.AddConfigPath(".")                                  // 当前目录
	viper.AddConfigPath("../config")                          // 上级目录的 config 目录
	viper.AddConfigPath("../../config")                       // 上上级目录的 config 目录

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var opts Options
	if err := viper.Unmarshal(&opts); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &opts, nil
}
