package main

import (
	"context"
	"fmt"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/config"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/connection"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/model"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/repository"
	"log"
	"os"
	"time"
)

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
	repo := repository.NewAccountPostgresRepository(dbConnectionData)

	_, err = dbConnectionData.BunDB.NewCreateTable().Model((*model.AccountModel)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		panic(err)
	}

	_, err = dbConnectionData.BunDB.NewCreateTable().Model((*model.TransactionModel)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		panic(err)
	}

	// Insert Account
	newAcc := &account.Account{DocumentNumber: "00000000000"}
	acc, err := repo.Save(ctx, newAcc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Account inserido: %+v\n", acc)

	acc, err = repo.FindByID(ctx, 1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Account Pesquisado: %+v\n", acc)
	repo2 := repository.NewTransactionPostgresRepository(dbConnectionData)
	newtransaction := &transaction.Transaction{
		AccountID:       1,
		OperationTypeID: transaction.Purchase,
		Amount:          100.00,
		EventDate:       time.Now(),
	}
	trs, err := repo2.Save(ctx, newtransaction)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Transaction inserido: %+v\n", trs)
}
