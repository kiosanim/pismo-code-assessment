package account

import (
	"context"
)

type AccountRepository interface {
	FindByID(ctx context.Context, accountID int64) (*Account, error)
	FindByDocumentNumber(ctx context.Context, documentNumber string) (*Account, error)
	Save(ctx context.Context, newAccount *Account) (*Account, error)
}
