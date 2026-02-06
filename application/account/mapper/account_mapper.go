package mapper

import (
	"github.com/kiosanim/pismo-code-assessment/application/account/dto"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
)

func CreateDTOToEntity(req dto.CreateAccountRequest) *account.Account {
	return &account.Account{
		DocumentNumber: req.DocumentNumber,
	}
}

func CreateEntityToResponse(entity *account.Account) *dto.CreateAccountResponse {
	return &dto.CreateAccountResponse{
		AccountID:      entity.AccountID,
		DocumentNumber: entity.DocumentNumber,
	}
}

func FindEntityToResponse(entity *account.Account) *dto.FindAccountByIdResponse {
	return &dto.FindAccountByIdResponse{
		AccountID:      entity.AccountID,
		DocumentNumber: entity.DocumentNumber,
	}
}

func ListAccountsToResponse(entities []account.Account, limit int64, cursor int64) *dto.ListAccountsResponse {
	var accountsDTO []dto.AccountDTO

	for _, entity := range entities {
		accountDTO := dto.AccountDTO{
			AccountID:      entity.AccountID,
			DocumentNumber: entity.DocumentNumber,
		}
		accountsDTO = append(accountsDTO, accountDTO)
	}
	return &dto.ListAccountsResponse{
		Accounts: accountsDTO,
		Limit:    limit,
		Cursor:   cursor,
	}
}
