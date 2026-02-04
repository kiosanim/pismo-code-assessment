package service

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/application/account/dto"
	"github.com/kiosanim/pismo-code-assessment/internal/core/cache"
	"github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/logger/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"testing"
)

type AccountServiceTestSuite struct {
	suite.Suite
	repository *account.AccountRepositoryMock
	cache      cache.CacheRepository
	service    *AccountService
	ctx        context.Context
	log        *mock.MockLogger
}

func (s *AccountServiceTestSuite) SetupTest() {
	s.log = mock.NewMockLogger()
	s.repository = account.NewAccountRepositoryMock()
	control := gomock.NewController(s.T())
	defer control.Finish()
	s.cache = cache.NewCacheRepositoryMock(control)
	s.service = NewAccountService(s.repository, s.cache, s.log)
	s.ctx = context.Background()
}

func (s *AccountServiceTestSuite) TestNewAccountService() {
	repo := account.NewAccountRepositoryMock()
	service := NewAccountService(repo, s.cache, s.log)
	if service == nil {
		s.Fail("account service should not be nil")
	}
}

func (s *AccountServiceTestSuite) TestCreateAccountSuccess() {
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
	output, err := s.service.Create(s.ctx, input)
	s.NoError(err, "create account should return no error")
	documentNumberOutput := account.SanitizeDocumentNumber(output.DocumentNumber)
	s.NoError(account.IsValidDocumentNumber(documentNumberOutput), "document number should be valid")
	s.Greater(output.AccountID, accountIDToCompare, "account ID should be greater than zero")
}

func (s *AccountServiceTestSuite) TestCreateAccountInvalidParameters() {
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
	_, err := s.service.Create(s.ctx, input)
	s.Error(err, "create create_account should return no error")
}

func (s *AccountServiceTestSuite) TestFindByIDSuccess() {
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
	output, err := s.service.FindByID(s.ctx, input)
	s.NoError(err, "find account by ID should return no error")
	documentNumber := account.SanitizeDocumentNumber(output.DocumentNumber)
	s.NoError(account.IsValidDocumentNumber(documentNumber), "document number should be valid")
	s.Greater(output.AccountID, accountIDToCompare, "account ID should be greater than zero")
}

func (s *AccountServiceTestSuite) TestFindByIDNotFound() {
	var accountID int64 = 1
	s.repository.On(
		"FindByID",
		s.ctx,
		accountID,
	).Return(
		nil,
		errors.AccountNotFoundError)
	input := dto.FindAccountByIdRequest{AccountID: accountID}
	output, err := s.service.FindByID(s.ctx, input)
	s.ErrorIs(err, errors.AccountNotFoundError, errors.AccountNotFoundError.Error())
	s.Nil(output, "find account by ID should return nil because no Account was found")
}

func (s *AccountServiceTestSuite) TestFindByIDInvalidParameters() {
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
	output, err := s.service.FindByID(context.Background(), input)
	s.ErrorIs(err, errors.InvalidParametersError, "find account by ID should return error")
	s.Nil(output, "find account by ID should return nil because no Account was found")
}

func TestCreateAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountServiceTestSuite))
}
