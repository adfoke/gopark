package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config 结构体定义了应用的配置项
type Config struct {
	AppName string `mapstructure:"appname"` // Use mapstructure tags for Viper
	Port    int    `mapstructure:"port"`
	Debug   bool   `mapstructure:"debug"`
	Redis   string `mapstructure:"redis"`
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"database"`
}

// LoadConfig 加载配置文件并返回 Config 结构体
func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigName("config") // 配置文件名称（不包含扩展名）
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath(path)     // 使用传入的路径

	viper.AutomaticEnv() // 读取环境变量

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// 可选：在这里添加配置验证逻辑

	return config, nil
}
