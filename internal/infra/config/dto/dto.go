package dto

type DatabaseConfig struct {
	Dsn string `yaml:"dsn"`
}
type AppConfig struct {
	Env string `mapstructure:"env"`
}

type Configuration struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
}
