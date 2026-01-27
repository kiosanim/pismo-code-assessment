package factory

import (
	accountSvc "github.com/kiosanim/pismo-code-assessment/application/account/service"
	transactionSvc "github.com/kiosanim/pismo-code-assessment/application/transaction/service"
	"github.com/kiosanim/pismo-code-assessment/interfaces/http/handler"
	"github.com/kiosanim/pismo-code-assessment/internal/core/config"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
)

type Factory interface {
	Configuration() *config.Configuration
	AccountRepository() account.AccountRepository
	AccountService() accountSvc.AccountService
	TransactionRepository() transaction.TransactionRepository
	TransactionService() transactionSvc.TransactionService
	AccountHandler() handler.AccountHandler
	TransactionHandler() handler.TransactionHandler
}
