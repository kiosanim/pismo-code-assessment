package transaction

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/application/transaction/dto"
)

type Service interface {
	Create(ctx context.Context, input dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error)
	FindByID(ctx context.Context, request dto.FindTransactionByIdRequest) (*dto.FindTransactionByIdResponse, error)
}
