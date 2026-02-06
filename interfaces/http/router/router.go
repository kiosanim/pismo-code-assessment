package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/kiosanim/pismo-code-assessment/docs"
	"github.com/kiosanim/pismo-code-assessment/interfaces/http/handler"
	"github.com/kiosanim/pismo-code-assessment/interfaces/http/middleware"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func SetupRouter(accountHandler handler.AccountHandler, transactionHandler handler.TransactionHandler, log logger.Logger) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.TraceMiddleware())
	router.Use(middleware.LoggerMiddleware(log))
	api := router.Group("")
	{
		api.POST("/accounts", accountHandler.CreateAccount)
		api.GET("/accounts/:account_id", accountHandler.GetAccountByID)
		api.POST("/transactions", transactionHandler.CreateTransaction)
		api.GET("/transactions/:transaction_id", transactionHandler.GetTransactionByID)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
