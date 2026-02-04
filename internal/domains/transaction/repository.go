package transaction

import (
	"context"
)

type TransactionRepository interface {
	FindOperationTypeByID(ctx context.Context, operationTypeID int) (*OperationType, error)
	FindTransactionByID(ctx context.Context, transactionID int64) (*Transaction, error)
	Save(ctx context.Context, newTransaction *Transaction) (*Transaction, error)
}
