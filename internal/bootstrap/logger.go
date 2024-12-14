package bootstrap

import (
	"FastGo/config"
	"FastGo/internal/global"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// setupLogger 初始化日志
func setupLogger() error {
	logConfig := global.Config.Log

	// 设置日志输出格式
	encoder := getEncoder()

	// 设置日志级别
	var level zapcore.Level
	switch logConfig.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 设置日志输出方式
	writeSyncer := getWriteSyncer(&logConfig)

	// 创建核心
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 创建 logger
	global.Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// getEncoder 设置日志格式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getWriteSyncer 设置日志输出方式
func getWriteSyncer(logConfig *config.LogConfig) zapcore.WriteSyncer {
	// 获取项目根目录
	rootDir := getRootDir()
	
	// 将相对路径转换为绝对路径
	logPath := filepath.Join(rootDir, logConfig.FilePath)
	
	// 打印实际的日志路径，用于调试
	fmt.Printf("日志路径: %s\n", logPath)
	
	// 确保日志目录存在
	if err := os.MkdirAll(logPath, 0755); err != nil {
		fmt.Printf("创建日志目录失败: %v\n", err)
		// 失败时使用临时目录
		logPath = filepath.Join(os.TempDir(), "logs")
		if err := os.MkdirAll(logPath, 0755); err != nil {
			fmt.Printf("创建临时日志目录也失败了: %v\n", err)
		}
	}

	// 确保 app 子目录存在
	appLogDir := filepath.Join(logPath, "app")
	if err := os.MkdirAll(appLogDir, 0755); err != nil {
		fmt.Printf("创建应用日志目录失败: %v\n", err)
	}

	// 创建自定义的写入器
	writer := &DailyRotateWriter{
		Dir:      appLogDir,
		filename: "",
		date:     time.Now().Format("2006-01-02"),
	}

	var syncers []zapcore.WriteSyncer

	// 根据配置决定是否输出到控制台
	if logConfig.Console {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}

	// 添加文件输出
	syncers = append(syncers, zapcore.AddSync(writer))

	return zapcore.NewMultiWriteSyncer(syncers...)
}

// getRootDir 获取项目根目录
func getRootDir() string {
	// 获取当前文件的路径
	_, b, _, _ := runtime.Caller(0)
	// 获取项目根目录（假设当前文件在 internal/bootstrap 目录下）
	return filepath.Dir(filepath.Dir(filepath.Dir(b)))
}

// DailyRotateWriter 实现按天切割的写入器
type DailyRotateWriter struct {
	Dir      string    // 日志目录
	filename string    // 当前日志文件名
	file     *os.File  // 当前打开的文件
	date     string    // 当前日期
}

// Write 实现 io.Writer 接口
func (w *DailyRotateWriter) Write(p []byte) (n int, err error) {
	// 获取当前日期
	currentDate := time.Now().Format("2006-01-02")
	
	// 如果文件未打开或日期变更，需要切换文件
	if w.file == nil || w.date != currentDate {
		// 关闭之前的文件
		if w.file != nil {
			w.file.Close()
		}
		
		// 更新日期和文件名
		w.date = currentDate
		w.filename = w.date + ".log"
		
		// 获取完整的文件路径
		fullPath := filepath.Join(w.Dir, w.filename)
		
		// 检查文件是否存在
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			// 文件不存在，创建新文件
			w.file, err = os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return 0, fmt.Errorf("创建日志文件失败: %v", err)
			}
		} else {
			// 文件存在，直接打开
			w.file, err = os.OpenFile(fullPath, os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return 0, fmt.Errorf("打开��志文件失败: %v", err)
			}
		}
	}
	
	// 写入日志内容
	return w.file.Write(p)
}

// Sync 实现 zapcore.WriteSyncer 接口
func (w *DailyRotateWriter) Sync() error {
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}

// Close 关闭文件
func (w *DailyRotateWriter) Close() error {
	if w.file != nil {
		err := w.file.Close()
		w.file = nil
		return err
	}
	return nil
}
