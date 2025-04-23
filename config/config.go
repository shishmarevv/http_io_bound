package config

import (
	"time"

	"github.com/spf13/viper"
)

type HTTPConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type TaskConfig struct {
	WorkerCount int `mapstructure:"worker_count"`
}

type Config struct {
	API      HTTPConfig `mapstructure:"http"`
	IOserver HTTPConfig `mapstructure:"ioserver"`
	Task     TaskConfig `mapstructure:"task"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")
	viper.SetConfigType("yaml")

	viper.SetDefault("http.port", "8080")
	viper.SetDefault("http.read_timeout", "5s")
	viper.SetDefault("http.write_timeout", "10m")
	viper.SetDefault("http.idle_timeout", "1m")
	viper.SetDefault("ioserver.port", "9090")
	viper.SetDefault("ioserver.read_timeout", "5s")
	viper.SetDefault("ioserver.write_timeout", "10m")
	viper.SetDefault("ioserver.idle_timeout", "1m")
	viper.SetDefault("task.worker_count", 5)
	viper.SetDefault("cors.allow_origins", []string{"*"})
	viper.SetDefault("cors.allow_methods", []string{"GET", "POST", "OPTIONS"})
	viper.SetDefault("cors.allow_headers", []string{"Content-Type", "Authorization"})

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
