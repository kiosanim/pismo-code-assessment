package service

import (
	"context"
	goerrors "errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/cache"
	"github.com/kiosanim/pismo-code-assessment/internal/core/factory"
	"github.com/kiosanim/pismo-code-assessment/internal/core/lock"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"

	"github.com/kiosanim/pismo-code-assessment/application/transaction/dto"
	"github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	infralog "github.com/kiosanim/pismo-code-assessment/internal/infra/logger/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TransactionServiceTestSuite struct {
	suite.Suite
	accountRepository     *account.AccountRepositoryMock
	transactionRepository *transaction.TransactionRepositoryMock
	cache                 cache.CacheRepository
	service               *TransactionService
	ctx                   context.Context
	log                   logger.Logger
	factory               factory.Factory
	locker                lock.DistributedLockManager
}

func (s *TransactionServiceTestSuite) SetupTest() {
	//s.accountRepository = account.NewAccountRepositoryMock()
	s.log = infralog.NewMockLogger()
	control := gomock.NewController(s.T())
	defer control.Finish()
	s.cache = cache.NewCacheRepositoryMock(control)
	s.transactionRepository = transaction.NewTransactionRepositoryMock()
	dlm := lock.NewDistributedLockManagerMock(control)
	s.locker = dlm
	ft := factory.NewFactoryMock(control)
	s.factory = ft
	ft.EXPECT().TransactionRepository().Return(s.transactionRepository).AnyTimes()
	ft.EXPECT().CacheRepository().Return(s.cache).AnyTimes()
	ft.EXPECT().DistributedLockManager().Return(s.locker).AnyTimes()
	ft.EXPECT().Log().Return(s.log).AnyTimes()
	s.service = NewTransactionService(s.factory)
	s.ctx = context.Background()
}

func (s *TransactionServiceTestSuite) TestNewTransactionService() {
	service := NewTransactionService(s.factory)
	s.NotNil(service, "transaction service should not be nil")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionSuccess_Purchase() {
	var accountID int64 = 1
	var transactionID int64 = 100
	operationTypeID := transaction.Purchase
	amount := 123.45

	// Mock account exists
	s.accountRepository.On("FindByID", s.ctx, accountID).Return(
		&account.Account{AccountID: accountID, DocumentNumber: "12345678900"},
		nil,
	)

	// Mock operation type exists
	s.transactionRepository.On("FindOperationTypeByID", s.ctx, operationTypeID).Return(
		&transaction.OperationType{OperationTypeID: int64(operationTypeID), Description: "PURCHASE"},
		nil,
	)

	// Mock transaction save - Purchase operations become negative
	s.transactionRepository.On("Save", s.ctx, mock.MatchedBy(func(tx *transaction.Transaction) bool {
		return tx.AccountID == accountID &&
			tx.OperationTypeID == operationTypeID &&
			tx.Amount == -amount
	})).Return(
		&transaction.Transaction{
			TransactionID:   transactionID,
			AccountID:       accountID,
			OperationTypeID: operationTypeID,
			Amount:          -amount,
		},
		nil,
	)

	input := dto.CreateTransactionRequest{
		AccountID:       accountID,
		OperationTypeID: operationTypeID,
		Amount:          amount,
	}

	output, err := s.service.Create(s.ctx, input)

	s.NoError(err, "create transaction should return no error")
	s.NotNil(output, "output should not be nil")
	s.Equal(transactionID, output.Transaction.TransactionID, "transaction ID should match")
	s.Equal(accountID, output.Transaction.AccountID, "account ID should match")
	s.Equal(operationTypeID, output.Transaction.OperationTypeID, "operation type ID should match")
	// Amount in response should remain positive for display
	s.Equal(amount, output.Transaction.Amount, "amount should be positive in response")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionSuccess_Payment() {
	var accountID int64 = 1
	var transactionID int64 = 100
	operationTypeID := transaction.Payment
	amount := 500.00

	// Mock account exists
	s.accountRepository.On("FindByID", s.ctx, accountID).Return(
		&account.Account{AccountID: accountID, DocumentNumber: "12345678900"},
		nil,
	)

	// Mock operation type exists
	s.transactionRepository.On("FindOperationTypeByID", s.ctx, operationTypeID).Return(
		&transaction.OperationType{OperationTypeID: int64(operationTypeID), Description: "PAYMENT"},
		nil,
	)

	// Mock transaction save - Payment remains positive (code 4 in reverseAmountSign)
	s.transactionRepository.On("Save", s.ctx, mock.MatchedBy(func(tx *transaction.Transaction) bool {
		return tx.AccountID == accountID &&
			tx.OperationTypeID == operationTypeID &&
			tx.Amount == amount
	})).Return(
		&transaction.Transaction{
			TransactionID:   transactionID,
			AccountID:       accountID,
			OperationTypeID: operationTypeID,
			Amount:          amount,
		},
		nil,
	)

	input := dto.CreateTransactionRequest{
		AccountID:       accountID,
		OperationTypeID: operationTypeID,
		Amount:          amount,
	}

	output, err := s.service.Create(s.ctx, input)

	s.NoError(err, "create transaction should return no error")
	s.NotNil(output, "output should not be nil")
	s.Equal(amount, output.Transaction.Amount, "amount should remain positive for payment")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionError_InvalidAccountID() {
	input := dto.CreateTransactionRequest{
		AccountID:       -1, // Invalid
		OperationTypeID: transaction.Purchase,
		Amount:          100.0,
	}

	output, err := s.service.Create(s.ctx, input)

	s.Error(err, "should return error for invalid account ID")
	s.ErrorIs(err, errors.TransactionInvalidAccountIDError)
	s.Nil(output, "output should be nil")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionError_ZeroAccountID() {
	input := dto.CreateTransactionRequest{
		AccountID:       0, // Invalid
		OperationTypeID: transaction.Purchase,
		Amount:          100.0,
	}

	output, err := s.service.Create(s.ctx, input)

	s.Error(err, "should return error for zero account ID")
	s.ErrorIs(err, errors.TransactionInvalidAccountIDError)
	s.Nil(output, "output should be nil")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionError_NegativeAmount() {
	input := dto.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: transaction.Purchase,
		Amount:          -100.0, // Invalid negative amount
	}

	output, err := s.service.Create(s.ctx, input)

	s.Error(err, "should return error for negative amount")
	s.ErrorIs(err, errors.TransactionInvalidAmountNegativeError)
	s.Nil(output, "output should be nil")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionError_ZeroAmount() {
	input := dto.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: transaction.Purchase,
		Amount:          0, // Invalid zero amount
	}

	output, err := s.service.Create(s.ctx, input)

	s.Error(err, "should return error for zero amount")
	s.ErrorIs(err, errors.TransactionInvalidAmountNegativeError)
	s.Nil(output, "output should be nil")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionError_InvalidOperationType() {
	operationTypeID := 999 // Invalid operation type

	// Mock operation type not found
	s.transactionRepository.On("FindOperationTypeByID", s.ctx, operationTypeID).Return(
		nil,
		errors.TransactionInvalidOperationTypeError,
	)

	input := dto.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: operationTypeID,
		Amount:          100.0,
	}

	output, err := s.service.Create(s.ctx, input)

	s.Error(err, "should return error for invalid operation type")
	s.ErrorIs(err, errors.TransactionInvalidOperationTypeError)
	s.Nil(output, "output should be nil")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionError_AccountNotFound() {
	var accountID int64 = 999
	operationTypeID := transaction.Purchase

	// Mock account not found
	s.accountRepository.On("FindByID", s.ctx, accountID).Return(
		nil,
		errors.AccountNotFoundError,
	)

	// Mock operation type exists
	s.transactionRepository.On("FindOperationTypeByID", s.ctx, operationTypeID).Return(
		&transaction.OperationType{OperationTypeID: int64(operationTypeID), Description: "PURCHASE"},
		nil,
	)

	input := dto.CreateTransactionRequest{
		AccountID:       accountID,
		OperationTypeID: operationTypeID,
		Amount:          100.0,
	}

	output, err := s.service.Create(s.ctx, input)

	s.Error(err, "should return error when account not found")
	s.ErrorIs(err, errors.AccountNotFoundError)
	s.Nil(output, "output should be nil")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionError_RepositorySaveError() {
	var accountID int64 = 1
	operationTypeID := transaction.Purchase
	amount := 100.0

	// Mock account exists
	s.accountRepository.On("FindByID", s.ctx, accountID).Return(
		&account.Account{AccountID: accountID, DocumentNumber: "12345678900"},
		nil,
	)

	// Mock operation type exists
	s.transactionRepository.On("FindOperationTypeByID", s.ctx, operationTypeID).Return(
		&transaction.OperationType{OperationTypeID: int64(operationTypeID), Description: "PURCHASE"},
		nil,
	)

	// Mock repository save error
	repositoryError := goerrors.New("database error")
	s.transactionRepository.On("Save", s.ctx, mock.MatchedBy(func(tx *transaction.Transaction) bool {
		return tx.AccountID == accountID &&
			tx.OperationTypeID == operationTypeID &&
			tx.Amount == -amount
	})).Return(nil, repositoryError)

	input := dto.CreateTransactionRequest{
		AccountID:       accountID,
		OperationTypeID: operationTypeID,
		Amount:          amount,
	}

	output, err := s.service.Create(s.ctx, input)

	s.Error(err, "should return error from repository")
	s.ErrorIs(err, repositoryError)
	s.Nil(output, "output should be nil")
}

func (s *TransactionServiceTestSuite) TestIsAValidOperationType_Valid() {
	operationTypeID := transaction.Purchase

	s.transactionRepository.On("FindOperationTypeByID", s.ctx, operationTypeID).Return(
		&transaction.OperationType{OperationTypeID: int64(operationTypeID), Description: "PURCHASE"},
		nil,
	)

	result := s.service.isAValidOperationType(s.ctx, operationTypeID)
	s.True(result, "should return true for valid operation type")
}

func (s *TransactionServiceTestSuite) TestIsAValidOperationType_Invalid_NotFound() {
	operationTypeID := 999

	s.transactionRepository.On("FindOperationTypeByID", s.ctx, operationTypeID).Return(
		nil,
		goerrors.New("not found"),
	)

	result := s.service.isAValidOperationType(s.ctx, operationTypeID)
	s.False(result, "should return false for non-existent operation type")
}

func (s *TransactionServiceTestSuite) TestIsAValidOperationType_Invalid_Zero() {
	operationTypeID := 0

	result := s.service.isAValidOperationType(s.ctx, operationTypeID)
	s.False(result, "should return false for zero operation type")
}

func (s *TransactionServiceTestSuite) TestIsAValidOperationType_Invalid_Negative() {
	operationTypeID := -1

	result := s.service.isAValidOperationType(s.ctx, operationTypeID)
	s.False(result, "should return false for negative operation type")
}

func (s *TransactionServiceTestSuite) TestIsAValidOperationType_ReturnsNil() {
	operationTypeID := 5

	s.transactionRepository.On("FindOperationTypeByID", s.ctx, operationTypeID).Return(
		nil,
		nil,
	)

	result := s.service.isAValidOperationType(s.ctx, operationTypeID)
	s.False(result, "should return false when output is nil")
}

func (s *TransactionServiceTestSuite) TestReverseAmountSign_Purchase() {
	tx := &transaction.Transaction{
		OperationTypeID: transaction.Purchase,
		Amount:          100.0,
	}

	result := s.service.reverseAmountSign(tx)
	s.Equal(-100.0, result, "purchase amount should be negative")
	s.Equal(-100.0, tx.Amount, "transaction amount should be modified")
}

func (s *TransactionServiceTestSuite) TestReverseAmountSign_InstallmentPurchase() {
	tx := &transaction.Transaction{
		OperationTypeID: transaction.InstallmentPurchase,
		Amount:          200.0,
	}

	result := s.service.reverseAmountSign(tx)
	s.Equal(-200.0, result, "installment purchase amount should be negative")
}

func (s *TransactionServiceTestSuite) TestReverseAmountSign_Withdrawal() {
	tx := &transaction.Transaction{
		OperationTypeID: transaction.Withdrawal,
		Amount:          50.0,
	}

	result := s.service.reverseAmountSign(tx)
	s.Equal(-50.0, result, "withdrawal amount should be negative")
}

func (s *TransactionServiceTestSuite) TestReverseAmountSign_Payment() {
	tx := &transaction.Transaction{
		OperationTypeID: transaction.Payment,
		Amount:          500.0,
	}

	result := s.service.reverseAmountSign(tx)
	s.Equal(500.0, result, "payment amount should remain positive")
	s.Equal(500.0, tx.Amount, "transaction amount should not be modified for payment")
}

func TestTransactionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceTestSuite))
}

func TestTransactionService_Create(t1 *testing.T) {
	type fields struct {
		accountRepository     account.AccountRepository
		transactionRepository transaction.TransactionRepository
		componentName         string
		log                   logger.Logger
	}
	type args struct {
		ctx     context.Context
		request dto.CreateTransactionRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dto.CreateTransactionResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TransactionService{
				accountRepository:     tt.fields.accountRepository,
				transactionRepository: tt.fields.transactionRepository,
				componentName:         tt.fields.componentName,
				log:                   tt.fields.log,
			}
			got, err := t.Create(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionService_FindByID(t1 *testing.T) {
	type fields struct {
		accountRepository     account.AccountRepository
		transactionRepository transaction.TransactionRepository
		componentName         string
		log                   logger.Logger
	}
	type args struct {
		ctx     context.Context
		request dto.FindTransactionByIdRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dto.FindTransactionByIdResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TransactionService{
				accountRepository:     tt.fields.accountRepository,
				transactionRepository: tt.fields.transactionRepository,
				componentName:         tt.fields.componentName,
				log:                   tt.fields.log,
			}
			got, err := t.FindByID(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t1.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("FindByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionService_isAValidOperationType(t1 *testing.T) {
	type fields struct {
		accountRepository     account.AccountRepository
		transactionRepository transaction.TransactionRepository
		componentName         string
		log                   logger.Logger
	}
	type args struct {
		ctx             context.Context
		operationTypeID int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TransactionService{
				accountRepository:     tt.fields.accountRepository,
				transactionRepository: tt.fields.transactionRepository,
				componentName:         tt.fields.componentName,
				log:                   tt.fields.log,
			}
			if got := t.isAValidOperationType(tt.args.ctx, tt.args.operationTypeID); got != tt.want {
				t1.Errorf("isAValidOperationType() = %v, want %v", got, tt.want)
			}
		})
	}
}
