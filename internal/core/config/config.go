package config

type DatabaseConfig struct {
	Dsn string `yaml:"dsn"`
}
type AppConfig struct {
	Env      string `mapstructure:"env"`
	Address  string `mapstructure:"address"`
	LogLevel string `mapstructure:"logLevel"`
}

type Configuration struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
}

type Config interface {
	LoadConfig(path string) (*Configuration, error)
	MustLoadConfig(path string) *Configuration
}
