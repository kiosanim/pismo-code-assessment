package account

import (
	"context"
)

type AccountRepository interface {
	FindByID(ctx context.Context, accountID int64) (*Account, error)
	Save(ctx context.Context, account *Account) (*Account, error)
}
