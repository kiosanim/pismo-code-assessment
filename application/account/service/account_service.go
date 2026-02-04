package service

import (
	"context"
	"errors"
	"github.com/kiosanim/pismo-code-assessment/application/account/dto"
	"github.com/kiosanim/pismo-code-assessment/application/account/mapper"
	"github.com/kiosanim/pismo-code-assessment/internal/core/cache"
	coreerr "github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
)

type AccountService struct {
	accountRepository account.AccountRepository
	cache             cache.CacheRepository
	componentName     string
	log               logger.Logger
}

func NewAccountService(repository account.AccountRepository, cache cache.CacheRepository, log logger.Logger) *AccountService {
	accountService := AccountService{accountRepository: repository, cache: cache, log: log}
	accountService.componentName = logger.ComponentNameFromStruct(accountService)
	return &accountService
}

func (a *AccountService) FindByID(ctx context.Context, request dto.FindAccountByIdRequest) (*dto.FindAccountByIdResponse, error) {
	a.log.Debug(a.componentName+".FindByID", "request", request)
	if request.AccountID <= 0 {
		return nil, coreerr.InvalidParametersError
	}
	output, err := a.accountRepository.FindByID(ctx, request.AccountID)
	if err != nil {
		return nil, err
	}
	if output == nil {
		return nil, coreerr.AccountNotFoundError
	}
	response := mapper.FindEntityToResponse(output)
	return response, nil
}

func (a *AccountService) Create(ctx context.Context, request dto.CreateAccountRequest) (*dto.CreateAccountResponse, error) {
	a.log.Debug(a.componentName+".Create", "request", request)
	if request.DocumentNumber == "" || account.IsValidDocumentNumber(request.DocumentNumber) != nil {
		return nil, coreerr.InvalidParametersError
	}
	accountByDocumentNumber, err := a.accountRepository.FindByDocumentNumber(ctx, request.DocumentNumber)
	if err != nil && !errors.Is(err, coreerr.AccountNotFoundError) {
		return nil, err
	}
	if accountByDocumentNumber != nil {
		return nil, coreerr.AccountAlreadyExistsForDocumentNumberError
	}
	accountRequest := mapper.CreateDTOToEntity(request)
	output, err := a.accountRepository.Save(ctx, accountRequest)
	if err != nil {
		return nil, err
	}
	response := mapper.CreateEntityToResponse(output)
	return response, nil
}
