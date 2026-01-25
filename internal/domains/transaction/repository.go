package transaction

import (
	"context"
	"errors"
)

var TransactionRepositoryInvalidParametersError = errors.New("transaction repository invalid parameters")

type TransactionRepository interface {
	Save(ctx context.Context, transaction *Transaction) (*Transaction, error)
}
