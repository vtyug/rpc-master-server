package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Options struct {
	MySQL    MySQLOptions `json:"mysql"`
	Redis    RedisOptions `json:"redis"`
	Log      LogConfig    `json:"log"`
	Jwt      JWT          `json:"jwt"`
	Language string       `json:"language"`
}

func New() (*Options, error) {
	// 获取程序运行时的工作目录
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("获取工作目录失败: %w", err)
	}

	viper.SetConfigName("config") // 配置文件名称
	viper.SetConfigType("yaml")   // 配置文件类型

	// 添加配置文件搜索路径，按优先级从高到低
	viper.AddConfigPath(filepath.Join(workDir, "config"))     // 工作目录下的 config
	viper.AddConfigPath("config")                             // 相对于执行程序的 config 目录
	viper.AddConfigPath(".")                                  // 当前目录

	// 通过环境变量指定配置文件路径
	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		viper.AddConfigPath(configPath)
	}

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
