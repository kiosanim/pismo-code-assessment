package service

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/application/transaction/dto"
	"github.com/kiosanim/pismo-code-assessment/application/transaction/mapper"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
)

type TransactionService struct {
	accountRepository     account.AccountRepository
	transactionRepository transaction.TransactionRepository
}

func NewTransactionService(accountRepository account.AccountRepository, transactionRepository transaction.TransactionRepository) *TransactionService {
	return &TransactionService{accountRepository: accountRepository, transactionRepository: transactionRepository}
}

func (t *TransactionService) Create(ctx context.Context, request dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error) {
	if request.AccountID <= 0 || request.OperationTypeID <= 0 || request.Amount <= 0.0 || !transaction.IsAValidTransactionType(request.OperationTypeID) {
		return nil, transaction.TransactionServiceInvalidParametersError
	}
	_, err := t.accountRepository.FindByID(ctx, request.AccountID)
	if err != nil {
		return nil, err
	}
	newTransaction := mapper.CreateDTOToEntity(request)
	if newTransaction.OperationTypeID != transaction.Payment {
		newTransaction.Amount *= -1
	}
	response, err := t.transactionRepository.Save(ctx, newTransaction)
	if err != nil {
		return nil, err
	}
	return mapper.EntityToResponse(response), nil
}
