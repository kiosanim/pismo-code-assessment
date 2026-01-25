package config

import (
	"fmt"
	"github.com/kiosanim/pismo-code-assessment/internal/core/config"
	"github.com/spf13/viper"
	"os"
)

func LoadConfig(path string) (*config.Configuration, error) {
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
	var cfg config.Configuration
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, config.ConfigFileUnmarshalError
	}
	return &cfg, nil
}

func MustLoadConfig(path string) *config.Configuration {
	cfg, err := LoadConfig(path)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	return cfg
}
