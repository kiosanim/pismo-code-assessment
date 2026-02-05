package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/kiosanim/pismo-code-assessment/interfaces/http/router"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/factory"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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
	appFactory := factory.NewAppFactory(ctx)
	sLogger := logger.NewSlogLogger(ctx, appFactory.Configuration())
	r := router.NewRouterFactory(appFactory, sLogger)
	server := &http.Server{
		Addr:    appFactory.Configuration().App.Address,
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
	sLogger.Warn("Closing Database Connection...")
	defer closeDBConnection(&appFactory, *sLogger)
	sLogger.Warn("Closing Cache Connection...")
	defer closeCacheConnection(&appFactory, *sLogger)
	err := server.Shutdown(ctxWithTimeout)
	if err != nil {
		sLogger.Error(fmt.Sprintf("Server Shutdown Error: %v", err))
	}
	sLogger.Warn("Server Shutdown Completed")
}

func closeDBConnection(appFactory *factory.AppFactory, log logger.SlogLogger) {
	err := appFactory.ConnectionData().Db.Close()
	if err != nil {
		log.Error(fmt.Sprintf("Failed to close Database Connection: %v", err))
	}
	return
}

func closeCacheConnection(appFactory *factory.AppFactory, log logger.SlogLogger) {
	err := appFactory.CacheConnectionData().Rdb.Close()
	if err != nil {
		log.Error(fmt.Sprintf("Failed to close Cache Connection: %v", err))
	}
	return
}
