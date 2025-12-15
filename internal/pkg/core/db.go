package core

import (
	"fmt"
	"time"

	"backend-go/pkg/config"
	"backend-go/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	cfg := config.C.MySQL

	// 自定义 GORM Logger，使用 Zap
	newLogger := gormlogger.New(
		ZapGormWriter{},
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormlogger.Info,
			IgnoreRecordNotFoundError: true, // 忽略 RecordNotFound 错误日志 (因为我们会手动处理)
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		logger.Log.Fatal("failed to connect database", zap.Error(err))
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second)

	// 注册 AuditPlugin 自动填充 Creator、Updater、TenantID
	if err := db.Use(&AuditPlugin{}); err != nil {
		logger.Log.Fatal("failed to register AuditPlugin", zap.Error(err))
	}

	DB = db
	logger.Info("Database connected successfully")
	return db
}

// ZapGormWriter 适配 GORM Logger 接口
type ZapGormWriter struct{}

func (w ZapGormWriter) Printf(message string, data ...interface{}) {
	logger.Info(fmt.Sprintf(message, data...))
}
