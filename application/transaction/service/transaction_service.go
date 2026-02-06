package service

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/application/transaction/dto"
	"github.com/kiosanim/pismo-code-assessment/application/transaction/mapper"
	"github.com/kiosanim/pismo-code-assessment/internal/core/cache"
	"github.com/kiosanim/pismo-code-assessment/internal/core/contextutils"
	coreerr "github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/factory"
	"github.com/kiosanim/pismo-code-assessment/internal/core/lock"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	"time"
)

const purchaseOperationCode = 4

type TransactionService struct {
	accountRepository     account.AccountRepository
	transactionRepository transaction.TransactionRepository
	cache                 cache.CacheRepository
	componentName         string
	locker                lock.DistributedLockManager
	log                   logger.Logger
}

func NewTransactionService(factory factory.Factory) *TransactionService {
	return &TransactionService{
		componentName:         "TransactionService",
		accountRepository:     factory.AccountRepository(),
		transactionRepository: factory.TransactionRepository(),
		cache:                 factory.CacheRepository(),
		locker:                factory.DistributedLockManager(),
		log:                   factory.Log(),
	}
}

func (t *TransactionService) Create(ctx context.Context, request dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error) {
	traceID := contextutils.GetTraceID(ctx)
	t.log.Debug(t.componentName+".Create", "request", request, "x_trace_id", traceID)
	err := t.validateRequestParameters(ctx, request)
	if err != nil {
		t.log.Warn(t.componentName+".Create", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	_, err = t.accountRepository.FindByID(ctx, request.AccountID)
	if err != nil {
		t.log.Warn(t.componentName+".Create", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	newTransaction := mapper.CreateDTOToEntity(request)
	newTransaction.Amount = t.reverseAmountSign(newTransaction) //Change the amount sign for debt operations
	newTransaction.EventDate = time.Now()
	lck, err := t.locker.WaitToLockUsingDefaultTimeConfiguration(ctx, lock.TransactionCreationLockKey)
	if err != nil {
		t.log.Warn(t.componentName+".Create", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	response, err := t.transactionRepository.Save(ctx, newTransaction)
	if err != nil {
		t.log.Warn(t.componentName+".Create", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	err = t.locker.Unlock(ctx, lck)
	if err != nil {
		t.log.Warn(t.componentName+".Create", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	response.Amount = t.reverseAmountSign(newTransaction) //Returning value sign only for user presentation
	return mapper.EntityToResponse(response), nil
}

func (t *TransactionService) isAValidOperationType(ctx context.Context, operationTypeID int) bool {
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

func (t *TransactionService) FindByID(ctx context.Context, request dto.FindTransactionByIdRequest) (*dto.FindTransactionByIdResponse, error) {
	traceID := contextutils.GetTraceID(ctx)
	t.log.Debug(t.componentName+".FindByID", "request", request, "x_trace_id", traceID)
	if request.TransactionID <= 0 {
		err := coreerr.InvalidParametersError
		t.log.Warn(t.componentName+".FindByID", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	output, err := t.transactionRepository.FindTransactionByID(ctx, request.TransactionID)
	if err != nil {
		t.log.Warn(t.componentName+".FindByID", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	if output == nil {
		err := coreerr.AccountNotFoundError
		t.log.Warn(t.componentName+".FindByID", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	response := mapper.EntityByIdToResponseById(output)
	return response, nil
}

func (t *TransactionService) validateRequestParameters(ctx context.Context, request dto.CreateTransactionRequest) error {
	if request.AccountID <= 0 {
		return coreerr.TransactionInvalidAccountIDError
	}
	if request.Amount <= 0 {
		return coreerr.TransactionInvalidAmountNegativeError
	}
	if !t.isAValidOperationType(ctx, request.OperationTypeID) {
		return coreerr.TransactionInvalidOperationTypeError
	}
	return nil
}

// reverseAmountSign Change the amount sign for debt operations
func (t *TransactionService) reverseAmountSign(newTransaction *transaction.Transaction) float64 {
	if newTransaction.OperationTypeID != purchaseOperationCode {
		newTransaction.Amount = -newTransaction.Amount
	}
	return newTransaction.Amount
}
