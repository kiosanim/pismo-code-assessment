package account

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/application/account/dto"
)

type Service interface {
	FindByID(ctx context.Context, input dto.FindAccountByIdRequest) (*dto.FindAccountByIdResponse, error)
	Create(ctx context.Context, input dto.CreateAccountRequest) (*dto.CreateAccountResponse, error)
}
