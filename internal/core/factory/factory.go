package factory

import (
	"github.com/kiosanim/pismo-code-assessment/interfaces/http/handler"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	"github.com/kiosanim/pismo-code-assessment/internal/core/cache"
	"github.com/kiosanim/pismo-code-assessment/internal/core/config"
	"github.com/kiosanim/pismo-code-assessment/internal/core/lock"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
)

type Factory interface {
	Configuration() *config.Configuration
	ConnectionData() *adapter.DatabaseConnectionData
	CacheConnectionData() *adapter.CacheConnectionData
	AccountRepository() account.AccountRepository
	TransactionRepository() transaction.TransactionRepository
	AccountHandler(accountService account.Service) *handler.AccountHandler
	TransactionHandler(transactionService transaction.Service) *handler.TransactionHandler
	CacheRepository() cache.CacheRepository
	DistributedLockManager() lock.DistributedLockManager
	Log() logger.Logger
}
