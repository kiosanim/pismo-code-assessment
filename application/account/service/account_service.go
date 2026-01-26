package service

import (
	"context"
	"errors"
	"github.com/kiosanim/pismo-code-assessment/application/account/dto"
	"github.com/kiosanim/pismo-code-assessment/application/account/mapper"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
)

type AccountService struct {
	accountRepository account.AccountRepository
	componentName     string
	log               logger.Logger
}

func NewAccountService(repository account.AccountRepository, log logger.Logger) *AccountService {
	accountService := &AccountService{accountRepository: repository, log: log}
	accountService.componentName = logger.ComponentNameFromStruct(accountService)
	return accountService
}

func (a *AccountService) FindByID(ctx context.Context, request dto.FindAccountByIdRequest) (*dto.FindAccountByIdResponse, error) {
	a.log.Debug(a.componentName+".FindByID", "request", request)
	if request.AccountID <= 0 {
		return nil, account.AccountServiceInvalidParametersError
	}
	output, err := a.accountRepository.FindByID(ctx, request.AccountID)
	if err != nil {
		return nil, err
	}
	if output == nil {
		return nil, account.AccountServiceNotFoundError
	}
	response := mapper.FindEntityToResponse(output)

	return response, nil
}

func (a *AccountService) Create(ctx context.Context, request dto.CreateAccountRequest) (*dto.CreateAccountResponse, error) {
	a.log.Debug(a.componentName+".Create", "request", request)
	if request.DocumentNumber == "" || account.IsValidDocumentNumber(request.DocumentNumber) != nil {
		return nil, account.AccountServiceInvalidParametersError
	}
	accountByDocumentNumber, err := a.accountRepository.FindByDocumentNumber(ctx, request.DocumentNumber)
	if err != nil && !errors.Is(err, account.AccountServiceNotFoundError) {
		return nil, err
	}
	if accountByDocumentNumber != nil {
		return nil, account.AccountServiceAlreadyExistsForDocumentNumberError
	}
	accountRequest := mapper.CreateDTOToEntity(request)
	output, err := a.accountRepository.Save(ctx, accountRequest)
	if err != nil {
		return nil, err
	}
	response := mapper.CreateEntityToResponse(output)
	return response, nil
}
