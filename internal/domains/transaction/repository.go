package transaction

import (
	"context"
)

type TransactionRepository interface {
	Save(ctx context.Context, transaction *Transaction) (*Transaction, error)
}
