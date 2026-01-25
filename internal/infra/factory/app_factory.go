package factory

import (
	"github.com/kiosanim/pismo-code-assessment/application/account/service"
	"github.com/kiosanim/pismo-code-assessment/internal/core/config"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/connection"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/repository"
)

type AppFactory struct {
	configuration *config.Config
	connection    *connection.PostgresConnection
}

func NewAppFactory(configuration *config.Config, connection *connection.PostgresConnection) *AppFactory {
	connection.ConnectionData{}
	return &AppFactory{configuration: configuration, connection: connection}
}

func (a *AppFactory) Config() *config.Config {
	return a.configuration
}

func (a *AppFactory) AccountRepository() account.AccountRepository {
	return repository.NewAccountPostgresRepository(a.connection)
}

func (a *AppFactory) AccountService() service.AccountService {
	//TODO implement me
	panic("implement me")
}

func (a *AppFactory) TransactionRepository() transaction.TransactionRepository {
	//TODO implement me
	panic("implement me")
}
