package repository

import (
	"context"
	"database/sql"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	coreerr "github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/mapper"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/model"
)

type TransactionPostgresRepository struct {
	connectionData *adapter.DatabaseConnectionData
	componentName  string
	log            logger.Logger
}

func NewTransactionPostgresRepository(connectionData *adapter.DatabaseConnectionData, log logger.Logger) *TransactionPostgresRepository {
	repository := &TransactionPostgresRepository{
		connectionData: connectionData,
		log:            log,
	}
	repository.componentName = logger.ComponentNameFromStruct(repository)
	return repository
}

func (t *TransactionPostgresRepository) Save(ctx context.Context, newTransaction *transaction.Transaction) (*transaction.Transaction, error) {
	t.log.Debug(t.componentName+".FindByID", "request", newTransaction)
	transactionModel := mapper.ToTransactionModel(newTransaction)
	if transactionModel == nil {
		return nil, coreerr.InvalidParametersError
	}

	tx, err := t.connectionData.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, coreerr.DatabaseCreateTransactionError
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "insert into transactions(account_id, operation_type, amount, event_date) values(?, ?, ?, ?)")
	if err != nil {
		return nil, coreerr.DatabasePrepareStatementError
	}
	defer stmt.Close()
	err = stmt.QueryRowContext(
		ctx,
		transactionModel.AccountID,
		transactionModel.OperationType,
		transactionModel.Amount,
		transactionModel.EventDate).Scan(
		&transactionModel.AccountID,
		&transactionModel.TransactionID,
		&transactionModel.OperationTypeID,
		&transactionModel.Amount,
		&transactionModel.EventDate)
	if err != nil {
		return nil, coreerr.DatabaseInsertionError
	}
	err = tx.Commit()
	if err != nil {
		return nil, coreerr.DatabaseFailToCommitError
	}
	return mapper.ToTransactionEntity(transactionModel), nil
}

func (t *TransactionPostgresRepository) FindOperationTypeByID(ctx context.Context, operationTypeID int) (*transaction.OperationType, error) {
	t.log.Debug(t.componentName+".FindByID", "request", operationTypeID)
	var selectedOperationType model.OperationTypeModel
	tx, err := t.connectionData.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, coreerr.DatabaseCreateTransactionError
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "select operation_type_id, description from operation_types where operation_type_id = ?")
	if err != nil {
		return nil, coreerr.DatabasePrepareStatementError
	}
	defer stmt.Close()
	err = stmt.QueryRowContext(ctx, operationTypeID).Scan(&selectedOperationType.OperationTypeID, &selectedOperationType.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, coreerr.OperationTypeNotFoundError
		}
		t.log.Debug(t.componentName+".FindByID", "error", err)
		return nil, coreerr.DatabaseQueryError
	}
	return mapper.ToOperationTypeEntity(&selectedOperationType), nil
}

func (t *TransactionPostgresRepository) FindTransactionByID(ctx context.Context, transactionID int64) (*transaction.Transaction, error) {
	t.log.Debug(t.componentName+".FindByID", "request", transactionID)
	var transactionModel model.TransactionModel
	tx, err := t.connectionData.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, coreerr.DatabaseCreateTransactionError
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "select transaction_id, account_id, operation_type, amount from transactions where transaction_id = ?")
	if err != nil {
		return nil, coreerr.DatabasePrepareStatementError
	}
	defer stmt.Close()
	err = stmt.QueryRowContext(
		ctx,
		transactionID).Scan(
		&transactionModel.TransactionID,
		&transactionModel.AccountID,
		&transactionModel.OperationTypeID,
		transactionModel.Amount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, coreerr.TransactionNotFoundError
		}
		t.log.Debug(t.componentName+".FindByID", "error", err)
		return nil, coreerr.DatabaseQueryError
	}
	return mapper.ToTransactionEntity(&transactionModel), nil
}
