package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/kiosanim/pismo-code-assessment/interfaces/http/router"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	corecfg "github.com/kiosanim/pismo-code-assessment/internal/core/config"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/config"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/connection"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/factory"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setupConfig() *corecfg.Configuration {
	path, _ := os.Getwd()
	cfg, err := config.LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func setupCache(ctx context.Context, cfg *corecfg.Configuration) *adapter.CacheConnectionData {
	var conn adapter.CacheConnection = connection.NewRedisConnection(cfg)
	cacheConnectionData, err := conn.Connect(ctx)
	if err != nil {
		panic(err)
	}
	return cacheConnectionData
}

func setupDatabase(cfg *corecfg.Configuration) *adapter.DatabaseConnectionData {
	var conn adapter.DatabaseConnection = connection.NewPostgresConnection(cfg)
	dbConnectionData, err := conn.Connect()
	if err != nil {
		panic(err)
	}
	return dbConnectionData
}

// @title           Pismo Code Assessment API
// @version         1.0
// @description     Customer Account & Transactions
// @termsOfService  http://swagger.io/terms/

// @contact.name   Fábio Sartori
// @contact.url    https://github.com/kiosanim

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /
func main() {
	ctx := context.Background()
	cfg := setupConfig()
	sLogger := logger.NewSlogLogger(ctx, cfg)
	dbConnectionData := setupDatabase(cfg)
	cacheConnectionData := setupCache(ctx, cfg)
	appFactory := factory.NewAppFactory(cfg, dbConnectionData, cacheConnectionData, sLogger)
	if appFactory == nil {
		sLogger.Error("App Factory not initialized")
		os.Exit(1)
	}
	r := router.NewRouterFactory(*appFactory, sLogger)
	server := &http.Server{
		Addr:    cfg.App.Address,
		Handler: r,
	}
	go func() {
		sLogger.Info(fmt.Sprintf("Server Listening on: %s", server.Addr))
		err := server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			sLogger.Error(fmt.Sprintf("ListenAndServe Error: %v", err))
			os.Exit(2)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sLogger.Warn("Shutdown Server...")
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sLogger.Info("Closing Database Connection...")
	defer dbConnectionData.Db.Close()
	sLogger.Info("Closing Database Connection...")
	defer cacheConnectionData.Rdb.Close()
	err := server.Shutdown(ctxWithTimeout)
	if err != nil {
		sLogger.Error(fmt.Sprintf("Server Shutdown Error: %v", err))
	}
	sLogger.Info("Server Shutdown Completed")
}
