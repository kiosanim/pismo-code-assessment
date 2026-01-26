package factory

import (
	accSvc "github.com/kiosanim/pismo-code-assessment/application/account/service"
	trnSvc "github.com/kiosanim/pismo-code-assessment/application/transaction/service"
	"github.com/kiosanim/pismo-code-assessment/interfaces/http/handler"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	"github.com/kiosanim/pismo-code-assessment/internal/core/config"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/repository"
)

type AppFactory struct {
	configuration  *config.Configuration
	connectionData *adapter.ConnectionData
	log            logger.Logger
}

func NewAppFactory(configuration *config.Configuration, connectionData *adapter.ConnectionData, log logger.Logger) *AppFactory {
	return &AppFactory{
		configuration:  configuration,
		connectionData: connectionData,
		log:            log,
	}
}

func (a *AppFactory) Configuration() *config.Configuration {
	return a.configuration
}

func (a *AppFactory) AccountRepository() account.AccountRepository {
	return repository.NewAccountPostgresRepository(
		a.connectionData,
		a.log,
	)
}

func (a *AppFactory) AccountService() *accSvc.AccountService {
	return accSvc.NewAccountService(
		a.AccountRepository(),
		a.log,
	)
}

func (a *AppFactory) TransactionRepository() transaction.TransactionRepository {
	return repository.NewTransactionPostgresRepository(
		a.connectionData,
		a.log,
	)
}

func (a *AppFactory) TransactionService() *trnSvc.TransactionService {
	return trnSvc.NewTransactionService(
		a.AccountRepository(),
		a.TransactionRepository(),
		a.log,
	)
}

func (a *AppFactory) AccountHandler() *handler.AccountHandler {
	return handler.NewAccountHandler(
		a.AccountService(),
		a.log,
	)
}

func (a *AppFactory) TransactionHandler() *handler.TransactionHandler {
	return handler.NewTransactionHandler(
		a.TransactionService(),
		a.log,
	)
}
