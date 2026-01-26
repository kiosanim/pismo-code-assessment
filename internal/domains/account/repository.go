package account

import (
	"context"
	"errors"
)

var AccountRepositoryInvalidParametersError = errors.New("account repository invalid parameters")

type AccountRepository interface {
	FindByID(ctx context.Context, accountID int64) (*Account, error)
	FindByDocumentNumber(ctx context.Context, documentNumber string) (*Account, error)
	Save(ctx context.Context, newAccount *Account) (*Account, error)
}
