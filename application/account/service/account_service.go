package service

import (
	"context"
	"errors"
	"github.com/kiosanim/pismo-code-assessment/application/account/dto"
	"github.com/kiosanim/pismo-code-assessment/application/account/mapper"
	"github.com/kiosanim/pismo-code-assessment/internal/core/cache"
	"github.com/kiosanim/pismo-code-assessment/internal/core/contextutils"
	coreerr "github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/factory"
	"github.com/kiosanim/pismo-code-assessment/internal/core/lock"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
)

type AccountService struct {
	accountRepository account.AccountRepository
	cache             cache.CacheRepository
	componentName     string
	locker            lock.DistributedLockManager
	log               logger.Logger
}

func NewAccountService(factory factory.Factory) *AccountService {
	return &AccountService{
		componentName:     "AccountService",
		accountRepository: factory.AccountRepository(),
		cache:             factory.CacheRepository(),
		locker:            factory.DistributedLockManager(),
		log:               factory.Log(),
	}
}

func (a *AccountService) FindByID(ctx context.Context, request dto.FindAccountByIdRequest) (*dto.FindAccountByIdResponse, error) {
	traceID := contextutils.GetTraceID(ctx)
	a.log.Debug(a.componentName+".FindByID", "request", request, "x_trace_id", traceID)
	if request.AccountID <= 0 {
		err := coreerr.InvalidParametersError
		a.log.Warn(a.componentName+".FindByID", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	output, err := a.accountRepository.FindByID(ctx, request.AccountID)
	if err != nil {
		a.log.Warn(a.componentName+".FindByID", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	if output == nil {
		err := coreerr.AccountNotFoundError
		a.log.Warn(a.componentName+".FindByID", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	response := mapper.FindEntityToResponse(output)
	return response, nil
}

func (a *AccountService) Create(ctx context.Context, request dto.CreateAccountRequest) (*dto.CreateAccountResponse, error) {
	traceID := contextutils.GetTraceID(ctx)
	a.log.Debug(a.componentName+".Create", "request", request, "x_trace_id", traceID)
	if request.DocumentNumber == "" || account.IsValidDocumentNumber(request.DocumentNumber) != nil {
		err := coreerr.InvalidParametersError
		a.log.Warn(a.componentName+".Create", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	accountByDocumentNumber, err := a.accountRepository.FindByDocumentNumber(ctx, request.DocumentNumber)
	if err != nil && !errors.Is(err, coreerr.AccountNotFoundError) {
		a.log.Warn(a.componentName+".Create", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	if accountByDocumentNumber != nil {
		err := coreerr.AccountAlreadyExistsForDocumentNumberError
		a.log.Warn(a.componentName+".Create", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	accountRequest := mapper.CreateDTOToEntity(request)
	lck, err := a.locker.WaitToLockUsingDefaultTimeConfiguration(ctx, lock.AccountCreationLockKey)
	if err != nil {
		a.log.Warn(a.componentName+".Create", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	output, err := a.accountRepository.Save(ctx, accountRequest)
	if err != nil {
		a.log.Warn(a.componentName+".Create", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	err = a.locker.Unlock(ctx, lck)
	if err != nil {
		a.log.Warn(a.componentName+".Create", "error", err, "x_trace_id", traceID)
		return nil, err
	}
	response := mapper.CreateEntityToResponse(output)
	return response, nil
}

func (a *AccountService) List(ctx context.Context, request dto.ListAccountsRequest) (*dto.ListAccountsResponse, error) {
	traceID := contextutils.GetTraceID(ctx)
	a.log.Debug(a.componentName+".List", "request", request, "x_trace_id", traceID)
	var lastID int64 = 0
	if request.Cursor <= 0 {
		return nil, coreerr.InvalidParametersError
	}
	accounts, err := a.accountRepository.List(ctx, request.Limit, lastID)
	if err != nil {
		return nil, err
	}
	nextCursor := int64(0)
	if len(accounts) > 0 {
		nextCursor = accounts[len(accounts)-1].AccountID
	}
	return mapper.ListAccountsToResponse(accounts, request.Limit, nextCursor), nil
}
