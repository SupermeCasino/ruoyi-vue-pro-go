package logger

import (
	"os"

	"github.com/wxlbd/ruoyi-mall-go/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

func Init() {
	cfg := config.C.Log

	// 文件写入器 (支持轮转)
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize, // MB
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge, // Days
		Compress:   true,
	})

	// 编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 核心配置
	var core zapcore.Core
	if config.C.App.Env == "local" {
		// 本地开发：同时输出到控制台和文件，使用 Console 编码
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), parseLevel(cfg.Level)),
			zapcore.NewCore(fileEncoder, writeSyncer, parseLevel(cfg.Level)),
		)
	} else {
		// 生产环境：只输出到文件，使用 JSON 编码
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			writeSyncer,
			parseLevel(cfg.Level),
		)
	}

	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	// 替换全局 Logger
	zap.ReplaceGlobals(Log)
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// 辅助函数
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

// ... 可以添加更多辅助函数
// NewLogger 提供 Logger 实例 (Wire Provider)
func NewLogger() *zap.Logger {
	return Log
}
