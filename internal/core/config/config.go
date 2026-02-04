package config

import "time"

type DistributedLock struct {
	TTL           time.Duration `mapstructure:"ttl_ms"`
	RetryInterval time.Duration `mapstructure:"retry_interval_ms"`
}

type CacheConfig struct {
	URL string `mapstructure:"url"`
}

type DatabaseConfig struct {
	URL string `yaml:"url"`
}
type AppConfig struct {
	Env      string `mapstructure:"env"`
	Address  string `mapstructure:"address"`
	LogLevel string `mapstructure:"log_level"`
}

type Configuration struct {
	App             AppConfig       `mapstructure:"app"`
	Database        DatabaseConfig  `mapstructure:"database"`
	Cache           CacheConfig     `mapstructure:"cache"`
	DistributedLock DistributedLock `mapstructure:"distributed_lock"`
}

type Config interface {
	LoadConfig(path string) (*Configuration, error)
	MustLoadConfig(path string) *Configuration
}
