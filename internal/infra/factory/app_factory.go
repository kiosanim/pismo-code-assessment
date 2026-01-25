package factory

import (
	accSvc "github.com/kiosanim/pismo-code-assessment/application/account/service"
	trnSvc "github.com/kiosanim/pismo-code-assessment/application/transaction/service"
	"github.com/kiosanim/pismo-code-assessment/interfaces/http/handler"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	"github.com/kiosanim/pismo-code-assessment/internal/core/config"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/repository"
)

type AppFactory struct {
	configuration  *config.Configuration
	connectionData *adapter.ConnectionData
}

func NewAppFactory(configuration *config.Configuration, connectionData *adapter.ConnectionData) *AppFactory {
	return &AppFactory{configuration: configuration, connectionData: connectionData}
}

func (a *AppFactory) AppConfiguration() *config.Configuration {
	return a.configuration
}

func (a *AppFactory) AccountRepository() account.AccountRepository {
	return repository.NewAccountPostgresRepository(a.connectionData)
}

func (a *AppFactory) AccountService() *accSvc.AccountService {
	return accSvc.NewAccountService(repository.NewAccountPostgresRepository(a.connectionData))
}

func (a *AppFactory) TransactionRepository() transaction.TransactionRepository {
	return repository.NewTransactionPostgresRepository(a.connectionData)
}

func (a *AppFactory) TransactionService() *trnSvc.TransactionService {
	return trnSvc.NewTransactionService(repository.NewAccountPostgresRepository(a.connectionData), repository.NewTransactionPostgresRepository(a.connectionData))
}

func (a *AppFactory) AccountHandler() *handler.AccountHandler {
	return handler.NewAccountHandler(a.AccountService())
}

func (a *AppFactory) TransactionHandler() *handler.TransactionHandler {
	return handler.NewTransactionHandler(a.TransactionService())
}
