package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kiosanim/pismo-code-assessment/interfaces/http/handler"
)

func SetupRouter(accountHandler handler.AccountHandler, transactionHandler handler.TransactionHandler) *gin.Engine {
	router := gin.Default()
	api := router.Group("")
	{
		api.POST("/accounts", accountHandler.CreateAccount)
		api.GET("/accounts/:account_id", accountHandler.GetAccountByID)
		api.POST("/transactions", transactionHandler.CreateTransaction)
	}
	return router
}
