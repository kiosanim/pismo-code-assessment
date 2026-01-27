package account

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/application/account/dto"
)

type Service interface {
	FindByID(ctx context.Context, request dto.FindAccountByIdRequest) (*dto.FindAccountByIdResponse, error)
	Create(ctx context.Context, response dto.CreateAccountRequest) (*dto.CreateAccountResponse, error)
}
