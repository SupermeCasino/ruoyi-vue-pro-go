package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var C = new(Config)

type Config struct {
	App   AppConfig   `mapstructure:"app"`
	HTTP  HTTPConfig  `mapstructure:"http"`
	Log   LogConfig   `mapstructure:"log"`
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
	Trade TradeConfig `mapstructure:"trade"`
	Pay   PayConfig   `mapstructure:"pay"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type HTTPConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	DSN         string `mapstructure:"dsn"`
	MaxIdle     int    `mapstructure:"max_idle"`
	MaxOpen     int    `mapstructure:"max_open"`
	MaxLifetime int    `mapstructure:"max_lifetime"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type TradeConfig struct {
	Express ExpressConfig `mapstructure:"express"`
}

type ExpressConfig struct {
	Client string      `mapstructure:"client"`
	Kd100  Kd100Config `mapstructure:"kd100"`
	// KdNiao KdNiaoConfig `mapstructure:"kdniao"`
}

type Kd100Config struct {
	Customer string `mapstructure:"customer"`
	Key      string `mapstructure:"key"`
}

type PayConfig struct {
	OrderNotifyURL  string `mapstructure:"order_notify_url"`
	RefundNotifyURL string `mapstructure:"refund_notify_url"`
	OrderNoPrefix   string `mapstructure:"order_no_prefix"`
	WalletPayAppKey string `mapstructure:"wallet_pay_app_key"`
}

func Load() error {
	// 读取环境变量
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local"
	}

	viper.SetConfigName("config." + env) // e.g. config.local
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")    // 相对路径
	viper.AddConfigPath("../config") // 兼容测试路径

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(C); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
