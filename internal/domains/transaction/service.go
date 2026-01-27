package transaction

import (
	"context"
	"errors"
	"github.com/kiosanim/pismo-code-assessment/application/transaction/dto"
)

var TransactionServiceInvalidParametersError = errors.New("invalid parameter")
var TransactionServiceInvalidOperationTypeError = errors.New("invalid operation type")
var TransactionServiceInvalidAmountNegativeError = errors.New("invalid amount. must be a positive value")
var TransactionServiceInvalidAccountIDError = errors.New("invalid account ID")

type Service interface {
	Create(ctx context.Context, input dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error)
}
