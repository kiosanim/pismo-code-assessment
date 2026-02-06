package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	"github.com/kiosanim/pismo-code-assessment/internal/core/contextutils"
	coreerr "github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/mapper"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/model"
)

type AccountPostgresRepository struct {
	connectionData *adapter.DatabaseConnectionData
	componentName  string
	log            logger.Logger
}

func NewAccountPostgresRepository(connectionData *adapter.DatabaseConnectionData, log logger.Logger) *AccountPostgresRepository {
	repository := &AccountPostgresRepository{
		connectionData: connectionData,
		componentName:  "AccountPostgresRepository",
		log:            log,
	}
	repository.componentName = logger.ComponentNameFromStruct(repository)
	return repository
}

func (a *AccountPostgresRepository) FindByID(ctx context.Context, accountID int64) (*account.Account, error) {
	traceID := contextutils.GetTraceID(ctx)
	a.log.Debug(a.componentName+".FindByID", "accountID", accountID, "x_trace_id", traceID)
	var selectedAccount model.AccountModel
	stmt, err := a.connectionData.Db.PrepareContext(ctx, "SELECT account_id, document_number FROM accounts WHERE account_id = $1")
	if err != nil {
		a.log.Warn(a.componentName+".FindByID", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabasePrepareStatementError
	}
	defer stmt.Close()
	acc, err, done := a.validateAccountError(err)
	if done {
		return acc, err
	}
	err = stmt.QueryRowContext(ctx, accountID).Scan(&selectedAccount.AccountID, &selectedAccount.DocumentNumber)
	if err != nil {
		a.log.Warn(a.componentName+".FindByID", "error", err, "x_trace_id", traceID)
		if err == sql.ErrNoRows {
			return nil, coreerr.AccountNotFoundError
		}
		return nil, coreerr.DatabaseQueryError
	}
	return mapper.ToAccountEntity(&selectedAccount), nil
}

func (a *AccountPostgresRepository) validateAccountError(err error) (*account.Account, error, bool) {
	if errors.Is(err, sql.ErrNoRows) {
		return nil, coreerr.AccountNotFoundError, true
	} else if err != nil {
		return nil, err, true
	}
	return nil, nil, false
}

func (a *AccountPostgresRepository) FindByDocumentNumber(ctx context.Context, documentNumber string) (*account.Account, error) {
	traceID := contextutils.GetTraceID(ctx)
	a.log.Debug(a.componentName+".FindByDocumentNumber", "documentNumber", documentNumber, "x_trace_id", traceID)
	var selectedAccount model.AccountModel
	stmt, err := a.connectionData.Db.PrepareContext(ctx, "SELECT account_id, document_number FROM accounts WHERE document_number = $1")
	if err != nil {
		a.log.Warn(a.componentName+".FindByDocumentNumber", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabasePrepareStatementError
	}
	defer stmt.Close()
	acc, err, done := a.validateAccountError(err)
	if done {
		return acc, err
	}
	err = stmt.QueryRowContext(ctx, documentNumber).Scan(&selectedAccount.AccountID, &selectedAccount.DocumentNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			a.log.Warn(a.componentName+".FindByDocumentNumber", "error", err, "x_trace_id", traceID)
			return nil, coreerr.AccountNotFoundError
		}
		a.log.Warn(a.componentName+".FindByDocumentNumber", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabaseQueryError
	}

	return mapper.ToAccountEntity(&selectedAccount), nil
}

func (a *AccountPostgresRepository) Save(ctx context.Context, newAccount *account.Account) (*account.Account, error) {
	traceID := contextutils.GetTraceID(ctx)
	a.log.Debug(a.componentName+".Save", "newAccount", newAccount, "x_trace_id", traceID)
	accountModel := mapper.ToAccountModel(newAccount)
	if accountModel == nil {
		return nil, coreerr.InvalidParametersError
	}

	tx, err := a.connectionData.Db.BeginTx(ctx, nil)
	if err != nil {
		a.log.Warn(a.componentName+".Save", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabaseCreateTransactionError
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO accounts (document_number) VALUES ($1) RETURNING account_id, document_number;")
	if err != nil {
		a.log.Warn(a.componentName+".Save", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabasePrepareStatementError
	}
	defer stmt.Close()
	err = stmt.QueryRowContext(
		ctx,
		accountModel.DocumentNumber).Scan(
		&accountModel.AccountID,
		&accountModel.DocumentNumber)
	if err != nil {
		a.log.Warn(a.componentName+".Save", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabaseInsertionError
	}
	err = tx.Commit()
	if err != nil {
		a.log.Warn(a.componentName+".Save", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabaseFailToCommitError
	}
	return mapper.ToAccountEntity(accountModel), nil
}

func (a *AccountPostgresRepository) List(ctx context.Context, limit int, cursorID int64) ([]account.Account, error) {

	traceID := contextutils.GetTraceID(ctx)
	a.log.Debug(a.componentName+".List", "limit", limit, "cursorID", cursorID, "x_trace_id", traceID)
	tx, err := a.connectionData.Db.BeginTx(ctx, nil)
	if err != nil {
		a.log.Warn(a.componentName+".List", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabaseCreateTransactionError
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "SELECT account_id, document_number FROM accounts WHERE account_id > $1 ORDER BY account_id LIMIT $2")
	if err != nil {
		a.log.Warn(a.componentName+".List", "error", err, "x_trace_id", traceID)
		return nil, coreerr.DatabasePrepareStatementError
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, cursorID, limit)
	if err != nil {
		a.log.Warn(a.componentName+".List", "error", err, "x_trace_id", traceID)
		if err == sql.ErrNoRows {
			return nil, coreerr.AccountNotFoundError
		}
		a.log.Debug(a.componentName+".List", "error", err)
		return nil, coreerr.DatabaseQueryError
	}
	defer rows.Close()
	var accounts []account.Account
	for rows.Next() {
		var account model.AccountModel
		err = rows.Scan(&account.AccountID, &account.DocumentNumber)
		if err != nil {
			a.log.Warn(a.componentName+".List", "error", err, "x_trace_id", traceID)
			return nil, err
		}
		accounts = append(accounts, *mapper.ToAccountEntity(&account))
	}
	return accounts, nil
}
