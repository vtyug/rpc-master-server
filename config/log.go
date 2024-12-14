package config

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level" json:"level" yaml:"level"`                // 日志级别
	FilePath   string `mapstructure:"filepath" json:"filepath" yaml:"filepath"`       // 日志目录路径
	MaxSize    int    `mapstructure:"maxSize" json:"maxSize" yaml:"maxSize"`          // 每个日志文件保存的最大尺寸 单位：M
	MaxBackups int    `mapstructure:"maxBackups" json:"maxBackups" yaml:"maxBackups"` // 日志文件最多保存多少个备份
	MaxAge     int    `mapstructure:"maxAge" json:"maxAge" yaml:"maxAge"`             // 文件最多保存多少天
	Compress   bool   `mapstructure:"compress" json:"compress" yaml:"compress"`       // 是否压缩
	Console    bool   `mapstructure:"console" json:"console" yaml:"console"`          // 是否输出到控制台
}
