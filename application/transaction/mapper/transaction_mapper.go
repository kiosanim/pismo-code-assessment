package mapper

import (
	"github.com/kiosanim/pismo-code-assessment/application/transaction/dto"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
)

func CreateDTOToEntity(req dto.CreateTransactionRequest) *transaction.Transaction {
	return &transaction.Transaction{
		AccountID:       req.AccountID,
		Amount:          req.Amount,
		OperationTypeID: req.OperationTypeID,
	}
}

func EntityToResponse(entity *transaction.Transaction) *dto.CreateTransactionResponse {
	return &dto.CreateTransactionResponse{
		TransactionID:   entity.TransactionID,
		AccountID:       entity.AccountID,
		OperationTypeID: entity.OperationTypeID,
		Amount:          entity.Amount,
	}
}
