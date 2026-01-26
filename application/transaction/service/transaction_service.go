package service

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/application/transaction/dto"
	"github.com/kiosanim/pismo-code-assessment/application/transaction/mapper"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
)

const TRANSACTION_SERVICE_NAME = "TransactionService"

type TransactionService struct {
	accountRepository     account.AccountRepository
	transactionRepository transaction.TransactionRepository
	componentName         string
	log                   logger.Logger
}

func NewTransactionService(accountRepository account.AccountRepository, transactionRepository transaction.TransactionRepository, log logger.Logger) *TransactionService {
	transactionService := &TransactionService{
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		log:                   log,
	}
	transactionService.componentName = logger.ComponentNameFromStruct(transactionService)
	return transactionService
}

func (t *TransactionService) Create(ctx context.Context, request dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error) {
	t.log.Debug(t.componentName+".Create", "request", request)
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
