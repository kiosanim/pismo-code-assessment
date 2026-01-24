package config

import (
	"github.com/kiosanim/pismo-code-assessment/internal/core/config"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/config/dto"
	"github.com/spf13/viper"
	"os"
)

func LoadConfig(path string) (*dto.Configuration, error) {
	v := viper.New()
	env := os.Getenv("ENV")
	if env == "production" {
		v.SetConfigName("config.production")
	} else {
		v.SetConfigName("config")
	}
	v.SetConfigType("yaml")
	v.AddConfigPath(path)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, config.ConfigFileNotFountError
	}
	var cfg dto.Configuration
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, config.ConfigFileUnmarshalError
	}
	return &cfg, nil
}

func MustLoadConfig(path string) *dto.Configuration {
	cfg, err := LoadConfig(path)
	if err != nil {
		panic(err)
	}
	return cfg
}
