package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	coreerr "github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/mapper"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/model"
)

type TransactionPostgresRepository struct {
	connectionData *adapter.ConnectionData
	componentName  string
	log            logger.Logger
}

func NewTransactionPostgresRepository(connectionData *adapter.ConnectionData, log logger.Logger) *TransactionPostgresRepository {
	repository := &TransactionPostgresRepository{
		connectionData: connectionData,
		log:            log,
	}
	repository.componentName = logger.ComponentNameFromStruct(repository)
	return repository
}

func (a *TransactionPostgresRepository) Save(ctx context.Context, newTransaction *transaction.Transaction) (*transaction.Transaction, error) {
	a.log.Debug(a.componentName+".FindByID", "request", newTransaction)
	transactionModel := mapper.ToTransactionModel(newTransaction)
	if transactionModel == nil {
		return nil, coreerr.InvalidParametersError
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

func (a *TransactionPostgresRepository) FindOperationTypeByID(ctx context.Context, OperationTypeID int) (*transaction.OperationType, error) {
	a.log.Debug(a.componentName+".FindByID", "request", OperationTypeID)
	var selectedOperationType model.OperationTypeModel
	err := a.connectionData.BunDB.NewSelect().Model(&selectedOperationType).Where("operation_type_id = ?", OperationTypeID).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, coreerr.TransactionInvalidOperationTypeError
	} else if err != nil {
		return nil, err
	}
	return mapper.ToOperationTypeEntity(&selectedOperationType), nil
}
