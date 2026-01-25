package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/mapper"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/model"
)

type AccountPostgresRepository struct {
	connectionData *adapter.ConnectionData
}

func NewAccountPostgresRepository(connectionData *adapter.ConnectionData) *AccountPostgresRepository {
	return &AccountPostgresRepository{
		connectionData: connectionData,
	}
}

func (a *AccountPostgresRepository) FindByID(ctx context.Context, accountID int64) (*account.Account, error) {
	var selectedAccount model.AccountModel
	err := a.connectionData.BunDB.NewSelect().Model(&selectedAccount).Where("account_id = ?", accountID).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, account.AccountServiceNotFoundError
	} else if err != nil {
		return nil, err
	}
	return mapper.ToAccountEntity(&selectedAccount), nil
}

func (a *AccountPostgresRepository) Save(ctx context.Context, newAccount *account.Account) (*account.Account, error) {
	accountModel := mapper.ToAccountModel(newAccount)
	if accountModel == nil {
		return nil, account.AccountRepositoryInvalidParametersError
	}
	tx, err := a.connectionData.BunDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	_, err = tx.NewInsert().Model(accountModel).Returning("*").Exec(ctx)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return mapper.ToAccountEntity(accountModel), nil
}
