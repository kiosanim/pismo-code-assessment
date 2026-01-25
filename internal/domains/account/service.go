package account

import (
	"context"
	"errors"
	"github.com/kiosanim/pismo-code-assessment/application/account/dto"
)

var (
	AccountServiceInvalidParametersError = errors.New("account service invalid parameter")
	AccountServiceNotFoundError          = errors.New("account not found")
)

type Service interface {
	FindByID(ctx context.Context, request dto.FindAccountByIdRequest) (*dto.FindAccountByIdResponse, error)
	Create(ctx context.Context, response dto.CreateAccountRequest) (*dto.CreateAccountResponse, error)
}
