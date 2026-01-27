package account

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/application/account/dto"
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

func (m *AccountRepositoryMock) FindByDocumentNumber(ctx context.Context, documentNumber string) (*Account, error) {
	args := m.Called(ctx, documentNumber)
	val := args.Get(0)
	p, ok := val.(*Account)
	if !ok {
		return nil, args.Error(1)
	}
	return p, nil
}

type AccountServiceMock struct {
	mock.Mock
}

func NewAccountServiceMock() *AccountServiceMock {
	return &AccountServiceMock{}
}

func (m *AccountServiceMock) FindByID(ctx context.Context, request dto.FindAccountByIdRequest) (*dto.FindAccountByIdResponse, error) {
	args := m.Called(ctx, request)
	val := args.Get(0)
	p, ok := val.(*dto.FindAccountByIdResponse)
	if !ok {
		return nil, args.Error(1)
	}
	return p, nil
}

func (m *AccountServiceMock) Create(ctx context.Context, response dto.CreateAccountRequest) (*dto.CreateAccountResponse, error) {
	args := m.Called(ctx, response)
	val := args.Get(0)
	p, ok := val.(*dto.CreateAccountResponse)
	if !ok {
		return nil, args.Error(1)
	}
	return p, nil
}
