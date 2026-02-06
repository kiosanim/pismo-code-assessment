package account

import (
	"context"
)

type AccountRepository interface {
	FindByID(ctx context.Context, accountID int64) (*Account, error)
	FindByDocumentNumber(ctx context.Context, documentNumber string) (*Account, error)
	Save(ctx context.Context, newAccount *Account) (*Account, error)
	List(ctx context.Context, limit int64, cursorID int64) ([]Account, error)
}
