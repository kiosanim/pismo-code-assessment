package account

import (
	"context"
	"errors"
)

var AccountRepositoryInvalidParametersError = errors.New("account repository invalid parameters")

type AccountRepository interface {
	FindByID(ctx context.Context, accountID int64) (*Account, error)
	Save(ctx context.Context, account *Account) (*Account, error)
}
