package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config defines application configuration
type Config struct {
	AppName  string `mapstructure:"appname"` // Use mapstructure tags for Viper
	Port     int    `mapstructure:"port"`
	Debug    bool   `mapstructure:"debug"`
	Redis    string `mapstructure:"redis"`
	Database struct {
		Type string `mapstructure:"type"` // Database type, e.g. sqlite
		Path string `mapstructure:"path"` // SQLite database file path
	} `mapstructure:"database"`
}

// LoadConfig reads configuration and returns a Config
func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigName("config") // Base name without extension
	viper.SetConfigType("yaml")   // Config file format
	viper.AddConfigPath(path)     // Search path

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Apply defaults
	if config.Database.Type == "" {
		config.Database.Type = "sqlite"
	}
	if config.Database.Path == "" {
		config.Database.Path = "./gopark.db"
	}

	return config, nil
}
