package main

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/interfaces/http/router"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/config"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/connection"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/factory"
	"log"
	"os"
)

//
//// 1️⃣ Connect to DB (infra)
//db := database.NewPostgresConnection("postgres://user:pass@localhost:5432/mydb?sslmode=disable")
//
//// 2️⃣ Instantiate repository and service (application)
//accountRepo := database.NewAccountRepository(db)
//accountService := service.NewAccountService(accountRepo)
//
//// 3️⃣ Create handler (HTTP)
//accountHandler := handlers.NewAccountHandler(accountService)
//
//// 4️⃣ Setup routes
//r := routes.SetupRouter(accountHandler)
//
//// 5️⃣ Run server
//log.Println("🚀 Server running at :8080")
//if err := r.Run(":8080"); err != nil {
//log.Fatal(err)
//}
//
//<-ctx.Done()

func main() {
	path, _ := os.Getwd()

	cfg, err := config.LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}
	var conn adapter.Connection = connection.NewPostgresConnection(cfg)
	ctx := context.Background()
	dbConnectionData, err := conn.Connect(ctx)
	if err != nil {
		panic(err)
	}
	appFactory := factory.NewAppFactory(cfg, dbConnectionData)
	if appFactory == nil {
		panic("App Factory not initialized")
	}
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
	r := router.SetupRouter(*accountHandler, *transactionHandler)
	
	// 5️⃣ Run server
	log.Println("🚀 Server running at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()

}
