package repository

import (
	"context"
	"database/sql"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	"github.com/kiosanim/pismo-code-assessment/internal/core/contextutils"
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
	traceID := contextutils.GetTraceID(ctx)
	t.log.Debug(t.componentName+".Save", "newTransaction", newTransaction, "x_trace_id", traceID)
	transactionModel := mapper.ToTransactionModel(newTransaction)
	if transactionModel == nil {
		err := coreerr.InvalidParametersError
		t.log.Warn(t.componentName+".Save", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	tx, err := t.connectionData.Db.BeginTx(ctx, nil)
	if err != nil {
		t.log.Warn(t.componentName+".Save", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabaseCreateTransactionError
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO transactions(account_id, operation_type_id, amount, event_date) VALUES($1, $2, $3, $4) RETURNING transaction_id, account_id, operation_type_id, amount, event_date")
	if err != nil {
		t.log.Warn(t.componentName+".Save", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabasePrepareStatementError
	}
	defer stmt.Close()
	err = stmt.QueryRowContext(
		ctx,
		transactionModel.AccountID,
		transactionModel.OperationTypeID,
		transactionModel.Amount,
		transactionModel.EventDate).Scan(
		&transactionModel.TransactionID,
		&transactionModel.AccountID,
		&transactionModel.OperationTypeID,
		&transactionModel.Amount,
		&transactionModel.EventDate)
	if err != nil {
		t.log.Warn(t.componentName+".Save", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabaseInsertionError
	}
	err = tx.Commit()
	if err != nil {
		t.log.Warn(t.componentName+".Save", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabaseFailToCommitError
	}
	return mapper.ToTransactionEntity(transactionModel), nil
}

func (t *TransactionPostgresRepository) FindOperationTypeByID(ctx context.Context, operationTypeID int) (*transaction.OperationType, error) {
	traceID := contextutils.GetTraceID(ctx)
	t.log.Debug(t.componentName+".FindOperationTypeByID", "operationTypeID", operationTypeID, "x_trace_id", traceID)
	var selectedOperationType model.OperationTypeModel
	tx, err := t.connectionData.Db.BeginTx(ctx, nil)
	if err != nil {
		t.log.Warn(t.componentName+".FindOperationTypeByID", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabaseCreateTransactionError
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT operation_type_id, description FROM operation_types WHERE operation_type_id = $1")
	if err != nil {
		t.log.Warn(t.componentName+".FindOperationTypeByID", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabasePrepareStatementError
	}
	defer stmt.Close()
	err = stmt.QueryRowContext(ctx, operationTypeID).Scan(&selectedOperationType.OperationTypeID, &selectedOperationType.Description)
	if err != nil {
		t.log.Warn(t.componentName+".FindOperationTypeByID", "error", err, "x_trace_id", traceID)
		if err == sql.ErrNoRows {
			return nil, coreerr.OperationTypeNotFoundError
		}
		t.log.Debug(t.componentName+".FindByID", "error", err)
		return nil, coreerr.DatabaseQueryError
	}
	return mapper.ToOperationTypeEntity(&selectedOperationType), nil
}

func (t *TransactionPostgresRepository) FindTransactionByID(ctx context.Context, transactionID int64) (*transaction.Transaction, error) {
	traceID := contextutils.GetTraceID(ctx)
	t.log.Debug(t.componentName+".FindTransactionByID", "transactionID", transactionID, "x_trace_id", traceID)
	var transactionModel model.TransactionModel
	tx, err := t.connectionData.Db.BeginTx(ctx, nil)
	if err != nil {
		t.log.Warn(t.componentName+".FindTransactionByID", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabaseCreateTransactionError
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT transaction_id, account_id, operation_type_id, amount FROM transactions WHERE transaction_id = $1")
	if err != nil {
		t.log.Warn(t.componentName+".FindTransactionByID", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabasePrepareStatementError
	}
	defer stmt.Close()
	err = stmt.QueryRowContext(
		ctx,
		transactionID).Scan(
		&transactionModel.TransactionID,
		&transactionModel.AccountID,
		&transactionModel.OperationTypeID,
		&transactionModel.Amount)
	if err != nil {
		t.log.Warn(t.componentName+".FindTransactionByID", "error", err, "x_trace_id", traceID)
		if err == sql.ErrNoRows {
			return nil, coreerr.TransactionNotFoundError
		}
		return nil, coreerr.DatabaseQueryError
	}
	return mapper.ToTransactionEntity(&transactionModel), nil
}
