package service

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/application/transaction/dto"
	"github.com/kiosanim/pismo-code-assessment/application/transaction/mapper"
	"github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	"time"
)

const (
	changeValueMultiplier = -1.0
	purchaseOperationCode = 4
)

type TransactionService struct {
	accountRepository     account.AccountRepository
	transactionRepository transaction.TransactionRepository
	componentName         string
	log                   logger.Logger
}

func NewTransactionService(accountRepository account.AccountRepository, transactionRepository transaction.TransactionRepository, log logger.Logger) *TransactionService {
	transactionService := TransactionService{
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		log:                   log,
	}
	transactionService.componentName = logger.ComponentNameFromStruct(transactionService)
	return &transactionService
}

func (t *TransactionService) Create(ctx context.Context, request dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error) {
	t.log.Debug(t.componentName+".Create", "request", request)
	err := t.validateRequestParameters(ctx, request)
	if err != nil {
		return nil, err
	}
	_, err = t.accountRepository.FindByID(ctx, request.AccountID)
	if err != nil {
		return nil, err
	}
	newTransaction := mapper.CreateDTOToEntity(request)
	newTransaction.Amount = t.reverseAmountSign(newTransaction) //Change the amount sign for debt operations
	newTransaction.EventDate = time.Now()
	response, err := t.transactionRepository.Save(ctx, newTransaction)
	if err != nil {
		return nil, err
	}
	response.Amount = t.reverseAmountSign(newTransaction) //Returning value sign only for user presentation
	return mapper.EntityToResponse(response), nil
}

func (t *TransactionService) isAValidOperationType(ctx context.Context, operationTypeID int) bool {
	t.log.Debug(t.componentName+".FindByID", "request", operationTypeID)
	if operationTypeID <= 0 {
		return false
	}
	output, err := t.transactionRepository.FindOperationTypeByID(ctx, operationTypeID)
	if err != nil {
		return false
	}
	if output == nil {
		return false
	}
	return true
}

func (t *TransactionService) validateRequestParameters(ctx context.Context, request dto.CreateTransactionRequest) error {
	if request.AccountID <= 0 {
		return errors.TransactionInvalidAccountIDError
	}
	if request.Amount <= 0 {
		return errors.TransactionInvalidAmountNegativeError
	}
	if !t.isAValidOperationType(ctx, request.OperationTypeID) {
		return errors.TransactionInvalidOperationTypeError
	}
	return nil
}

// reverseAmountSign Change the amount sign for debt operations
func (t *TransactionService) reverseAmountSign(newTransaction *transaction.Transaction) float64 {
	if newTransaction.OperationTypeID != purchaseOperationCode {
		newTransaction.Amount *= changeValueMultiplier
	}
	return newTransaction.Amount
}
