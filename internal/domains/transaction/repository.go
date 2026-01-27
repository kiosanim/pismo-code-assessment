package transaction

import (
	"context"
)

type TransactionRepository interface {
	FindOperationTypeByID(ctx context.Context, operationTypeID int) (*OperationType, error)
	Save(ctx context.Context, newTransaction *Transaction) (*Transaction, error)
}
