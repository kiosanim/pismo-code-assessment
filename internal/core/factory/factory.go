package factory

import (
	"github.com/kiosanim/pismo-code-assessment/application/account/service"
	"github.com/kiosanim/pismo-code-assessment/internal/core/config"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
)

type Factory interface {
	Config() *config.Config
	AccountRepository() account.AccountRepository
	AccountService() service.AccountService
	TransactionRepository() transaction.TransactionRepository
	//TransactionService() service.TransactionService
}
