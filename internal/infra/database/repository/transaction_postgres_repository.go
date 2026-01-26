package repository

import (
	"context"
	"database/sql"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/mapper"
)

const TRANSACTION_REPO_NAME = "TransactionRepository"

type TransactionPostgresRepository struct {
	connectionData *adapter.ConnectionData
	log            logger.Logger
}

func NewTransactionPostgresRepository(connectionData *adapter.ConnectionData, log logger.Logger) *TransactionPostgresRepository {
	return &TransactionPostgresRepository{
		connectionData: connectionData,
		log:            log,
	}
}

func (a *TransactionPostgresRepository) Save(ctx context.Context, newTransaction *transaction.Transaction) (*transaction.Transaction, error) {
	a.log.Debug(TRANSACTION_REPO_NAME+".Save", "request", newTransaction)
	transactionModel := mapper.ToTransactionModel(newTransaction)
	if transactionModel == nil {
		return nil, transaction.TransactionRepositoryInvalidParametersError
	}
	tx, err := a.connectionData.BunDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	_, err = tx.NewInsert().Model(transactionModel).Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return mapper.ToTransactionEntity(transactionModel), nil
}
