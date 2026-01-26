package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/factory"
)

const (
	ComponentName = "RouterFactory"
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
	accountSvc := appFactory.AccountService()
	if accountSvc == nil {
		panic("Account Service not initialized")
	}
	transactionSvc := appFactory.TransactionService()
	if transactionSvc == nil {
		panic("Transaction Service not initialized")
	}
	accountHandler := appFactory.AccountHandler()
	if accountHandler == nil {
		panic("Account Handler not initialized")
	}
	transactionHandler := appFactory.TransactionHandler()
	if transactionHandler == nil {
		panic("Transaction Handler not initialized")
	}
	return SetupRouter(*accountHandler, *transactionHandler, log)
}
