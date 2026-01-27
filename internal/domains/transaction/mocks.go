package transaction

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/application/transaction/dto"
	"github.com/stretchr/testify/mock"
)

type TransactionRepositoryMock struct {
	mock.Mock
}

func NewTransactionRepositoryMock() *TransactionRepositoryMock {
	return &TransactionRepositoryMock{}
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

func (tr *TransactionRepositoryMock) FindOperationTypeByID(ctx context.Context, operationTypeID int) (*OperationType, error) {
	args := tr.Called(ctx, operationTypeID)
	val := args.Get(0)
	p, ok := val.(*OperationType)
	if !ok {
		return nil, args.Error(1)
	}
	return p, nil
}

type TransactionServiceMock struct {
	mock.Mock
}

func NewTransactionServiceMock() *TransactionServiceMock {
	return &TransactionServiceMock{}
}

func (m *TransactionServiceMock) Create(ctx context.Context, input dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error) {
	args := m.Called(ctx, input)
	val := args.Get(0)
	p, ok := val.(*dto.CreateTransactionResponse)
	if !ok {
		return nil, args.Error(1)
	}
	return p, nil
}
