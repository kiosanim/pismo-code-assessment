package transaction

import (
	"context"
	"errors"
	"github.com/kiosanim/pismo-code-assessment/application/transaction/dto"
)

var TransactionServiceInvalidParametersError = errors.New("transaction service invalid parameter")

type Service interface {
	Create(ctx context.Context, input dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error)
}
