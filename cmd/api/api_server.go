package main

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/interfaces/http/router"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	corecfg "github.com/kiosanim/pismo-code-assessment/internal/core/config"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/config"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/connection"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/factory"
	"log"
	"os"
)

var ctx = context.Background()

func setupConfig() *corecfg.Configuration {
	path, _ := os.Getwd()
	cfg, err := config.LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func setupDatabase(cfg *corecfg.Configuration) *adapter.ConnectionData {
	var conn adapter.Connection = connection.NewPostgresConnection(cfg)
	dbConnectionData, err := conn.Connect(ctx)
	if err != nil {
		panic(err)
	}
	return dbConnectionData
}

func main() {
	cfg := setupConfig()
	dbConnectionData := setupDatabase(cfg)
	appFactory := factory.NewAppFactory(cfg, dbConnectionData)
	if appFactory == nil {
		panic("App Factory not initialized")
	}
	r := router.NewRouterFactory(*appFactory)

	log.Println("Server running at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
	<-ctx.Done()
}
