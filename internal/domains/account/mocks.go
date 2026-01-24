package account

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type AccountRepositoryMock struct {
	mock.Mock
}

func NewAccountRepositoryMock() *AccountRepositoryMock {
	return &AccountRepositoryMock{}
}

func (m *AccountRepositoryMock) FindByID(ctx context.Context, accountID int64) (*Account, error) {
	args := m.Called(ctx, accountID)
	val := args.Get(0)
	p, ok := val.(*Account)
	if !ok {
		return nil, args.Error(1)
	}
	return p, nil
}

func (m *AccountRepositoryMock) Save(ctx context.Context, account *Account) (*Account, error) {
	args := m.Called(ctx, account)
	val := args.Get(0)
	p, ok := val.(*Account)
	if !ok {
		return nil, args.Error(1)
	}
	return p, nil
}
