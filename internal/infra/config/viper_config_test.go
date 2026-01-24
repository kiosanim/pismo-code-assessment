package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	path := os.TempDir()
	configPath := filepath.Join(path, "config.yaml")
	configFileData := `
app:
  env: "development"
database:
  dsn: "postgres://user:password@localhost:5432/dbname?sslmode=disable"
`
	err := os.WriteFile(configPath, []byte(configFileData), 0644)
	if err != nil {
		t.Errorf("Failed to write test config file: %v", err)
	}
	config, err := LoadConfig(path)
	if err != nil {
		t.Errorf("Error loading config: %v", err)
	}
	tests := []struct {
		name    string
		got     any
		want    any
		wantErr bool
	}{
		{"must found variable", config.App.Env, "development", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.got; got != tt.want {
				t.Errorf("%s - want: %v - received: %v", tt.name, tt.want, got)
			}
		})
	}
}

func TestMustLoadConfig(t *testing.T) {
	path := os.TempDir()
	configPath := filepath.Join(path, "config.yaml")
	configFileData := `
app:
  env: "development"
database:
  dsn: "postgres://user:password@localhost:5432/dbname?sslmode=disable"
`
	err := os.WriteFile(configPath, []byte(configFileData), 0644)
	if err != nil {
		t.Errorf("Failed to write test config file: %v", err)
	}
	config := MustLoadConfig(path)
	tests := []struct {
		name    string
		got     any
		want    any
		wantErr bool
	}{
		{"must found variable", config.App.Env, "development", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.got; got != tt.want {
				t.Errorf("%s - want: %v - received: %v", tt.name, tt.want, got)
			}
		})
	}
}
