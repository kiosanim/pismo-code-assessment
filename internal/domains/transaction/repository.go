package transaction

import (
	"context"
	"errors"
)

var (
	OperationTypeRepositoryNotFoundError        = errors.New("operation type not found")
	TransactionRepositoryInvalidParametersError = errors.New("invalid parameters")
)

type TransactionRepository interface {
	FindOperationTypeByID(ctx context.Context, operationTypeID int) (*OperationType, error)
	Save(ctx context.Context, newTransaction *Transaction) (*Transaction, error)
}
