package config

import (
	"errors"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/config/dto"
)

var (
	ConfigFileNotFountError  = errors.New("config file not found")
	ConfigFileUnmarshalError = errors.New("config unmarshal error")
)

type Config interface {
	LoadConfig(path string) (*dto.Configuration, error)
	MustLoadConfig(path string) *dto.Configuration
}
