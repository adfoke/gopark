package config

import (
    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
)

type Config struct {
    AppName string
    Port    int
    Debug   bool
}

var AppConfig Config
var log = logrus.New()

func init() {
	// 配置 logrus
    log.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
    })
    viper.SetConfigName("config") // 配置文件名称（不包含扩展名）
    viper.SetConfigType("yaml")   // 配置文件类型
    viper.AddConfigPath("config") // 配置文件路径

    if err := viper.ReadInConfig(); err != nil {
        log.Fatalf("Error reading config file: %s", err)
    }

    if err := viper.Unmarshal(&AppConfig); err != nil {
        log.Fatalf("Unable to decode into struct: %s", err)
    }

    log.WithFields(logrus.Fields{
        "app_name": AppConfig.AppName,
        "port":     AppConfig.Port,
        "debug":    AppConfig.Debug,
    }).Infof("Config initialized")
}