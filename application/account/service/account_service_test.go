package service

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/application/account/dto"
	"github.com/kiosanim/pismo-code-assessment/internal/core/cache"
	"github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/factory"
	"github.com/kiosanim/pismo-code-assessment/internal/core/lock"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"testing"
)

type AccountServiceTestSuite struct {
	suite.Suite
	repository *account.AccountRepositoryMock
	cache      cache.CacheRepository
	ctx        context.Context
	log        *logger.LoggerMock
	factory    *factory.FactoryMock
	locker     *lock.DistributedLockManagerMock
}

func (s *AccountServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.ctx = context.Background()

	s.repository = account.NewAccountRepositoryMock()
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
	s.factory.EXPECT().AccountRepository().Return(s.repository).AnyTimes()
	s.factory.EXPECT().CacheRepository().Return(s.cache).AnyTimes()
	s.factory.EXPECT().DistributedLockManager().Return(s.locker).AnyTimes()
	s.factory.EXPECT().Log().Return(s.log).AnyTimes()
}

func (s *AccountServiceTestSuite) TestNewAccountService() {
	service := NewAccountService(s.factory)
	if service == nil {
		s.Fail("account service should not be nil")
	}
}

func (s *AccountServiceTestSuite) TestCreateAccountSuccess() {
	as := NewAccountService(s.factory)
	var accountID int64 = 1
	var accountIDToCompare int64 = 0
	documentNumber := "11987408098"
	s.repository.On(
		"Save",
		s.ctx,
		&account.Account{
			DocumentNumber: documentNumber,
		}).Return(
		&account.Account{
			AccountID:      accountID,
			DocumentNumber: documentNumber,
		}, nil)

	s.repository.On("FindByDocumentNumber", s.ctx, documentNumber).Return(
		nil, errors.AccountNotFoundError,
	)
	input := dto.CreateAccountRequest{DocumentNumber: documentNumber}
	output, err := as.Create(s.ctx, input)
	s.NoError(err, "create account should return no error")
	documentNumberOutput := account.SanitizeDocumentNumber(output.DocumentNumber)
	s.NoError(account.IsValidDocumentNumber(documentNumberOutput), "document number should be valid")
	s.Greater(output.AccountID, accountIDToCompare, "account ID should be greater than zero")
}

func (s *AccountServiceTestSuite) TestCreateAccountInvalidParameters() {
	service := NewAccountService(s.factory)
	var accountID int64 = 0
	s.repository.On(
		"Save",
		s.ctx,
		&account.Account{
			DocumentNumber: "",
		}).Return(
		&account.Account{
			AccountID:      accountID,
			DocumentNumber: "",
		}, errors.InvalidParametersError)
	input := dto.CreateAccountRequest{DocumentNumber: ""}
	_, err := service.Create(s.ctx, input)
	s.Error(err, "create create_account should return no error")
}

func (s *AccountServiceTestSuite) TestFindByIDSuccess() {
	service := NewAccountService(s.factory)
	var accountID int64 = 1
	var accountIDToCompare int64 = 0
	documentNumberFromMockResponse := "11987408098"
	s.repository.On(
		"FindByID",
		s.ctx,
		accountID).Return(
		&account.Account{
			AccountID:      accountID,
			DocumentNumber: documentNumberFromMockResponse,
		}, nil)
	input := dto.FindAccountByIdRequest{AccountID: accountID}
	output, err := service.FindByID(s.ctx, input)
	s.NoError(err, "find account by ID should return no error")
	documentNumber := account.SanitizeDocumentNumber(output.DocumentNumber)
	s.NoError(account.IsValidDocumentNumber(documentNumber), "document number should be valid")
	s.Greater(output.AccountID, accountIDToCompare, "account ID should be greater than zero")
}

func (s *AccountServiceTestSuite) TestFindByIDNotFound() {
	service := NewAccountService(s.factory)
	var accountID int64 = 1
	s.repository.On(
		"FindByID",
		s.ctx,
		accountID,
	).Return(
		nil,
		errors.AccountNotFoundError)
	input := dto.FindAccountByIdRequest{AccountID: accountID}
	output, err := service.FindByID(s.ctx, input)
	s.ErrorIs(err, errors.AccountNotFoundError, errors.AccountNotFoundError.Error())
	s.Nil(output, "find account by ID should return nil because no Account was found")
}

func (s *AccountServiceTestSuite) TestFindByIDInvalidParameters() {
	service := NewAccountService(s.factory)
	var accountID int64 = -1
	s.repository.On(
		"FindByID",
		context.Background(),
		&account.Account{
			AccountID: accountID,
		}).Return(
		nil,
		errors.InvalidParametersError)
	input := dto.FindAccountByIdRequest{AccountID: accountID}
	output, err := service.FindByID(context.Background(), input)
	s.ErrorIs(err, errors.InvalidParametersError, "find account by ID should return error")
	s.Nil(output, "find account by ID should return nil because no Account was found")
}

func TestCreateAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountServiceTestSuite))
}
