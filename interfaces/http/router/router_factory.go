package router

import (
	"github.com/gin-gonic/gin"
	acc "github.com/kiosanim/pismo-code-assessment/application/account/service"
	tra "github.com/kiosanim/pismo-code-assessment/application/transaction/service"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/factory"
)

func NewRouterFactory(appFactory factory.AppFactory, log logger.Logger) *gin.Engine {
	accountRepo := appFactory.AccountRepository()
	if accountRepo == nil {
		panic("Account Repository not initialized")
	}
	transactionRepo := appFactory.TransactionRepository()
	if transactionRepo == nil {
		panic("Transaction Repository not initialized")
	}
	accountSvc := acc.NewAccountService(&appFactory)
	if accountSvc == nil {
		panic("Account Service not initialized")
	}
	transactionSvc := tra.NewTransactionService(&appFactory)
	if transactionSvc == nil {
		panic("Transaction Service not initialized")
	}
	accountHandler := appFactory.AccountHandler(accountSvc)
	if accountHandler == nil {
		panic("Account Handler not initialized")
	}
	transactionHandler := appFactory.TransactionHandler(transactionSvc)
	if transactionHandler == nil {
		panic("Transaction Handler not initialized")
	}
	return SetupRouter(*accountHandler, *transactionHandler, log)
}
