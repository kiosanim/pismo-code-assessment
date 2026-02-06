package factory

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/interfaces/http/handler"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	"github.com/kiosanim/pismo-code-assessment/internal/core/cache"
	"github.com/kiosanim/pismo-code-assessment/internal/core/config"
	"github.com/kiosanim/pismo-code-assessment/internal/core/lock"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	infraconfig "github.com/kiosanim/pismo-code-assessment/internal/infra/config"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/connection"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/repository"
	infralock "github.com/kiosanim/pismo-code-assessment/internal/infra/lock"
	infralogger "github.com/kiosanim/pismo-code-assessment/internal/infra/logger"
	"log"
	"os"
)

type AppFactory struct {
	configuration       *config.Configuration
	connectionData      *adapter.DatabaseConnectionData
	cacheConnectionData *adapter.CacheConnectionData
	log                 logger.Logger
}

func NewAppFactory(ctx context.Context) AppFactory {
	appFactory := AppFactory{}
	configuration := appFactory.setupConfiguration()
	sLogger := infralogger.NewSlogLogger(ctx, configuration)
	connectionData := appFactory.setupDatabase(configuration)
	cacheConnectionData := appFactory.setupCache(ctx, configuration)
	appFactory.configuration = configuration
	appFactory.log = sLogger
	appFactory.connectionData = connectionData
	appFactory.cacheConnectionData = cacheConnectionData
	return appFactory
}

func (a *AppFactory) Configuration() *config.Configuration {
	return a.configuration
}

func (a *AppFactory) ConnectionData() *adapter.DatabaseConnectionData {
	return a.connectionData
}

func (a *AppFactory) CacheConnectionData() *adapter.CacheConnectionData {
	return a.cacheConnectionData
}

func (a *AppFactory) AccountRepository() account.AccountRepository {
	return repository.NewAccountPostgresRepository(
		a.connectionData,
		a.log,
	)
}

//
//func (a *AppFactory) AccountService() *accSvc.AccountService {
//	return accSvc.NewAccountService(
//		a.AccountRepository(),
//		a.CacheRepository(),
//		a.log,
//	)
//}

func (a *AppFactory) TransactionRepository() transaction.TransactionRepository {
	return repository.NewTransactionPostgresRepository(
		a.connectionData,
		a.log,
	)
}

//func (a *AppFactory) TransactionService() *trnSvc.TransactionService {
//	return trnSvc.NewTransactionService(
//		a.AccountRepository(),
//		a.TransactionRepository(),
//		a.CacheRepository(),
//		a.log,
//	)
//}

func (a *AppFactory) AccountHandler(accountService account.Service) *handler.AccountHandler {
	return handler.NewAccountHandler(
		accountService,
		a.log,
	)
}

func (a *AppFactory) TransactionHandler(transactionService transaction.Service) *handler.TransactionHandler {
	return handler.NewTransactionHandler(
		transactionService,
		a.log,
	)
}

func (a *AppFactory) CacheRepository() cache.CacheRepository {
	return repository.NewRedisRepository(
		a.cacheConnectionData,
		a.log,
	)
}

func (a *AppFactory) DistributedLockManager() lock.DistributedLockManager {
	return infralock.NewRedisDistributedLockManager(a.cacheConnectionData, a.configuration, a.log)
}

func (a *AppFactory) Log() logger.Logger {
	return a.log
}

func (a *AppFactory) setupCache(ctx context.Context, cfg *config.Configuration) *adapter.CacheConnectionData {
	var conn adapter.CacheConnection = connection.NewRedisConnection(cfg)
	cacheConnectionData, err := conn.Connect(ctx)
	if err != nil {
		panic(err)
	}
	return cacheConnectionData
}

func (a *AppFactory) setupDatabase(cfg *config.Configuration) *adapter.DatabaseConnectionData {
	var conn adapter.DatabaseConnection = connection.NewPostgresConnection(cfg)
	dbConnectionData, err := conn.Connect()
	if err != nil {
		panic(err)
	}
	return dbConnectionData
}

func (a *AppFactory) setupConfiguration() *config.Configuration {
	path, _ := os.Getwd()
	cfg, err := infraconfig.LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
