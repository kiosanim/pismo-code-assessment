package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
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
		log:            log,
	}
	repository.componentName = logger.ComponentNameFromStruct(repository)
	return repository
}

func (a *AccountPostgresRepository) FindByID(ctx context.Context, accountID int64) (*account.Account, error) {
	a.log.Debug(a.componentName+".FindByID", "request", accountID)
	var selectedAccount model.AccountModel
	stmt, err := a.connectionData.Db.PrepareContext(ctx, "select account_id, accounts.document_number from accounts where account_id = $1")
	if err != nil {
		return nil, coreerr.DatabasePrepareStatementError
	}
	defer stmt.Close()
	acc, err, done := a.validateAccountError(err)
	if done {
		return acc, err
	}
	err = stmt.QueryRowContext(ctx, accountID).Scan(&selectedAccount.AccountID, &selectedAccount.DocumentNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, coreerr.AccountNotFoundError
		}
		a.log.Error(a.componentName+".FindByID", "error", err)
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
	a.log.Debug(a.componentName+".FindByDocumentNumber", "request", documentNumber)
	var selectedAccount model.AccountModel
	stmt, err := a.connectionData.Db.PrepareContext(ctx, "select account_id, accounts.document_number from accounts where document_number = $1")
	if err != nil {
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
			return nil, coreerr.AccountNotFoundError
		}
		a.log.Error(a.componentName+".FindByDocumentNumberQueryRowContext", "err", err)
		return nil, coreerr.DatabaseQueryError
	}

	return mapper.ToAccountEntity(&selectedAccount), nil
}

func (a *AccountPostgresRepository) Save(ctx context.Context, newAccount *account.Account) (*account.Account, error) {
	a.log.Debug(a.componentName+".Save", "request", newAccount)
	accountModel := mapper.ToAccountModel(newAccount)
	if accountModel == nil {
		return nil, coreerr.InvalidParametersError
	}

	tx, err := a.connectionData.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, coreerr.DatabaseCreateTransactionError
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "insert into accounts (account_id, document_number) values ($1, $2)")
	if err != nil {
		return nil, coreerr.DatabasePrepareStatementError
	}
	defer stmt.Close()
	err = stmt.QueryRowContext(
		ctx,
		accountModel.AccountID,
		accountModel.DocumentNumber).Scan(
		&accountModel.AccountID,
		&accountModel.DocumentNumber)
	if err != nil {
		return nil, coreerr.DatabaseInsertionError
	}
	err = tx.Commit()
	if err != nil {
		return nil, coreerr.DatabaseFailToCommitError
	}
	return mapper.ToAccountEntity(accountModel), nil
}
