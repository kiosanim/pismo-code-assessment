package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const ConfigFileName = "config.yaml"

var FileContent = []byte(`app:
  env: "development"
  address: ":8080"
  loglevel: "debug"

database:
  dsn: "postgres://pismo:123456@postgres:5432/pismo_db?sslmode=disable"
`)

func main() {
	absoluteFilePath := filepath.Join(filepath.Dir(os.TempDir()), ConfigFileName)
	fmt.Printf("Creating a config.yaml file in a TEMP folder: %s\n", absoluteFilePath)
	err := os.WriteFile(absoluteFilePath, FileContent, os.FileMode(0644))
	if err != nil {
		fmt.Printf("Failed to create a config.yaml file at: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created config.yaml file at: %s\n", absoluteFilePath)
	fmt.Printf("Copy de file from: %s to the root of this project.\n", absoluteFilePath)
}
