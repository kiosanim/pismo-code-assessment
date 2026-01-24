package transaction

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type TransactionRepositoryMock struct {
	mock.Mock
}

func (tr *TransactionRepositoryMock) Save(ctx context.Context, transaction *Transaction) (*Transaction, error) {
	args := tr.Called(ctx, transaction)
	val := args.Get(0)
	p, ok := val.(*Transaction)
	if !ok {
		return nil, args.Error(1)
	}
	return p, nil
}
