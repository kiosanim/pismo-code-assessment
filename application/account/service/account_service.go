package service

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/application/account/dto"
	"github.com/kiosanim/pismo-code-assessment/application/account/mapper"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
)

type AccountService struct {
	repository account.AccountRepository
}

func NewAccountService(repository account.AccountRepository) *AccountService {
	return &AccountService{repository: repository}
}

func (a *AccountService) FindByID(ctx context.Context, request dto.FindAccountByIdRequest) (*dto.FindAccountByIdResponse, error) {
	if request.AccountID <= 0 {
		return nil, account.AccountServiceInvalidParametersError
	}
	output, err := a.repository.FindByID(ctx, request.AccountID)
	if err != nil {
		return nil, err
	}
	if output == nil {
		return nil, account.AccountNotFoundError
	}
	response := mapper.FindEntityToResponse(output)
	return response, nil
}

func (a *AccountService) Create(ctx context.Context, request dto.CreateAccountRequest) (*dto.CreateAccountResponse, error) {
	if request.DocumentNumber == "" || account.IsValidDocumentNumber(request.DocumentNumber) != nil {
		return nil, account.AccountServiceInvalidParametersError
	}
	accountRequest := mapper.CreateDTOToEntity(request)
	output, err := a.repository.Save(ctx, accountRequest)
	if err != nil {
		return nil, err
	}
	response := mapper.CreateEntityToResponse(output)
	return response, nil
}
