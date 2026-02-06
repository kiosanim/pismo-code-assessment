package config

type DistributedLock struct {
	TTL           int64 `mapstructure:"ttl_ms"`
	RetryInterval int64 `mapstructure:"retry_interval_ms"`
	WaitingTime   int64 `mapstructure:"waiting_time_ms"`
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
