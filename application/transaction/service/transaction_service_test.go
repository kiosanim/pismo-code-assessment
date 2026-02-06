package service

import (
	"context"
	"errors"
	"github.com/kiosanim/pismo-code-assessment/application/transaction/dto"
	"github.com/kiosanim/pismo-code-assessment/internal/core/cache"
	tranerr "github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/factory"
	"github.com/kiosanim/pismo-code-assessment/internal/core/lock"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"testing"
)

type TransactionServiceTestSuite struct {
	suite.Suite
	accountRepository     *account.AccountRepositoryMock
	transactionRepository *transaction.TransactionRepositoryMock
	cache                 cache.CacheRepository
	ctx                   context.Context
	log                   *logger.LoggerMock
	factory               *factory.FactoryMock
	locker                *lock.DistributedLockManagerMock
}

func (s *TransactionServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.ctx = context.Background()
	s.accountRepository = account.NewAccountRepositoryMock()
	s.transactionRepository = transaction.NewTransactionRepositoryMock()
	s.cache = cache.NewCacheRepositoryMock(ctrl)
	s.log = logger.NewLoggerMock(ctrl)
	s.locker = lock.NewDistributedLockManagerMock(ctrl)
	// Allow any number of these calls
	s.log.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	s.log.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	s.log.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
	s.locker.EXPECT().WaitToLockUsingDefaultTimeConfiguration(gomock.Any(), gomock.Any()).AnyTimes()
	s.locker.EXPECT().Unlock(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Factory returns same mocks
	s.factory = factory.NewFactoryMock(ctrl)
	s.factory.EXPECT().TransactionRepository().Return(s.transactionRepository).AnyTimes()
	s.factory.EXPECT().AccountRepository().Return(s.accountRepository).AnyTimes()
	s.factory.EXPECT().CacheRepository().Return(s.cache).AnyTimes()
	s.factory.EXPECT().DistributedLockManager().Return(s.locker).AnyTimes()
	s.factory.EXPECT().Log().Return(s.log).AnyTimes()
}

func (s *TransactionServiceTestSuite) TestNewTransactionService() {
	service := NewTransactionService(s.factory)
	s.NotNil(service, "transaction service should not be nil")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionSuccess_Purchase() {
	service := NewTransactionService(s.factory)
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
	output, err := service.Create(s.ctx, input)
	s.NoError(err, "create transaction should return no error")
	s.NotNil(output, "output should not be nil")
	s.Equal(transactionID, output.Transaction.TransactionID, "transaction ID should match")
	s.Equal(accountID, output.Transaction.AccountID, "account ID should match")
	s.Equal(operationTypeID, output.Transaction.OperationTypeID, "operation type ID should match")
	// Amount in response should remain positive for display
	s.Equal(amount, output.Transaction.Amount, "amount should be positive in response")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionSuccess_Payment() {
	service := NewTransactionService(s.factory)
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

	output, err := service.Create(s.ctx, input)

	s.NoError(err, "create transaction should return no error")
	s.NotNil(output, "output should not be nil")
	s.Equal(amount, output.Transaction.Amount, "amount should remain positive for payment")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionError_InvalidAccountID() {
	service := NewTransactionService(s.factory)
	input := dto.CreateTransactionRequest{
		AccountID:       -1, // Invalid
		OperationTypeID: transaction.Purchase,
		Amount:          100.0,
	}
	output, err := service.Create(s.ctx, input)
	s.Error(err, "should return error for invalid account ID")
	s.ErrorIs(err, tranerr.TransactionInvalidAccountIDError)
	s.Nil(output, "output should be nil")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionError_ZeroAccountID() {
	service := NewTransactionService(s.factory)
	input := dto.CreateTransactionRequest{
		AccountID:       0, // Invalid
		OperationTypeID: transaction.Purchase,
		Amount:          100.0,
	}
	output, err := service.Create(s.ctx, input)
	s.Error(err, "should return error for zero account ID")
	s.ErrorIs(err, tranerr.TransactionInvalidAccountIDError)
	s.Nil(output, "output should be nil")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionError_NegativeAmount() {
	service := NewTransactionService(s.factory)
	input := dto.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: transaction.Purchase,
		Amount:          -100.0, // Invalid negative amount
	}

	output, err := service.Create(s.ctx, input)

	s.Error(err, "should return error for negative amount")
	s.ErrorIs(err, tranerr.TransactionInvalidAmountNegativeError)
	s.Nil(output, "output should be nil")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionError_ZeroAmount() {
	service := NewTransactionService(s.factory)
	input := dto.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: transaction.Purchase,
		Amount:          0, // Invalid zero amount
	}

	output, err := service.Create(s.ctx, input)

	s.Error(err, "should return error for zero amount")
	s.ErrorIs(err, tranerr.TransactionInvalidAmountNegativeError)
	s.Nil(output, "output should be nil")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionError_InvalidOperationType() {
	service := NewTransactionService(s.factory)
	operationTypeID := 999 // Invalid operation type
	// Mock operation type not found
	s.transactionRepository.On("FindOperationTypeByID", s.ctx, operationTypeID).Return(
		nil,
		tranerr.TransactionInvalidOperationTypeError,
	)
	input := dto.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: operationTypeID,
		Amount:          100.0,
	}
	output, err := service.Create(s.ctx, input)
	s.Error(err, "should return error for invalid operation type")
	s.ErrorIs(err, tranerr.TransactionInvalidOperationTypeError)
	s.Nil(output, "output should be nil")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionError_AccountNotFound() {
	service := NewTransactionService(s.factory)
	var accountID int64 = 999
	operationTypeID := transaction.Purchase
	// Mock account not found
	s.accountRepository.On("FindByID", s.ctx, accountID).Return(
		nil,
		tranerr.AccountNotFoundError,
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
	output, err := service.Create(s.ctx, input)
	s.Error(err, "should return error when account not found")
	s.ErrorIs(err, tranerr.AccountNotFoundError)
	s.Nil(output, "output should be nil")
}

func (s *TransactionServiceTestSuite) TestCreateTransactionError_RepositorySaveError() {
	service := NewTransactionService(s.factory)
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
	repositoryError := errors.New("database error")
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
	output, err := service.Create(s.ctx, input)
	s.Error(err, "should return error from repository")
	s.ErrorIs(err, repositoryError)
	s.Nil(output, "output should be nil")
}

func (s *TransactionServiceTestSuite) TestIsAValidOperationType_Valid() {
	service := NewTransactionService(s.factory)
	operationTypeID := transaction.Purchase
	s.transactionRepository.On("FindOperationTypeByID", s.ctx, operationTypeID).Return(
		&transaction.OperationType{OperationTypeID: int64(operationTypeID), Description: "PURCHASE"},
		nil,
	)
	result := service.isAValidOperationType(s.ctx, operationTypeID)
	s.True(result, "should return true for valid operation type")
}

func (s *TransactionServiceTestSuite) TestIsAValidOperationType_Invalid_NotFound() {
	service := NewTransactionService(s.factory)
	operationTypeID := 999
	s.transactionRepository.On("FindOperationTypeByID", s.ctx, operationTypeID).Return(
		nil,
		errors.New("not found"),
	)
	result := service.isAValidOperationType(s.ctx, operationTypeID)
	s.False(result, "should return false for non-existent operation type")
}

func (s *TransactionServiceTestSuite) TestIsAValidOperationType_Invalid_Zero() {
	service := NewTransactionService(s.factory)
	operationTypeID := 0
	result := service.isAValidOperationType(s.ctx, operationTypeID)
	s.False(result, "should return false for zero operation type")
}

func (s *TransactionServiceTestSuite) TestIsAValidOperationType_Invalid_Negative() {
	service := NewTransactionService(s.factory)
	operationTypeID := -1
	result := service.isAValidOperationType(s.ctx, operationTypeID)
	s.False(result, "should return false for negative operation type")
}

func (s *TransactionServiceTestSuite) TestIsAValidOperationType_ReturnsNil() {
	service := NewTransactionService(s.factory)
	operationTypeID := 5

	s.transactionRepository.On("FindOperationTypeByID", s.ctx, operationTypeID).Return(
		nil,
		nil,
	)
	result := service.isAValidOperationType(s.ctx, operationTypeID)
	s.False(result, "should return false when output is nil")
}

func (s *TransactionServiceTestSuite) TestReverseAmountSign_Purchase() {
	service := NewTransactionService(s.factory)
	tx := &transaction.Transaction{
		OperationTypeID: transaction.Purchase,
		Amount:          100.0,
	}
	result := service.reverseAmountSign(tx)
	s.Equal(-100.0, result, "purchase amount should be negative")
	s.Equal(-100.0, tx.Amount, "transaction amount should be modified")
}

func (s *TransactionServiceTestSuite) TestReverseAmountSign_InstallmentPurchase() {
	service := NewTransactionService(s.factory)
	tx := &transaction.Transaction{
		OperationTypeID: transaction.InstallmentPurchase,
		Amount:          200.0,
	}
	result := service.reverseAmountSign(tx)
	s.Equal(-200.0, result, "installment purchase amount should be negative")
}

func (s *TransactionServiceTestSuite) TestReverseAmountSign_Withdrawal() {
	service := NewTransactionService(s.factory)
	tx := &transaction.Transaction{
		OperationTypeID: transaction.Withdrawal,
		Amount:          50.0,
	}
	result := service.reverseAmountSign(tx)
	s.Equal(-50.0, result, "withdrawal amount should be negative")
}

func (s *TransactionServiceTestSuite) TestReverseAmountSign_Payment() {
	service := NewTransactionService(s.factory)
	tx := &transaction.Transaction{
		OperationTypeID: transaction.Payment,
		Amount:          500.0,
	}
	result := service.reverseAmountSign(tx)
	s.Equal(500.0, result, "payment amount should remain positive")
	s.Equal(500.0, tx.Amount, "transaction amount should not be modified for payment")
}

func TestTransactionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceTestSuite))
}
