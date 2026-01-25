package mapper

import (
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/model"
)

func ToTransactionModel(entity *transaction.Transaction) *model.TransactionModel {
	if entity == nil {
		return nil
	}
	return &model.TransactionModel{
		AccountID:       entity.AccountID,
		TransactionID:   entity.TransactionID,
		OperationTypeID: entity.OperationTypeID,
		Amount:          entity.Amount,
		EventDate:       entity.EventDate,
	}
}

func ToTransactionEntity(model *model.TransactionModel) *transaction.Transaction {
	if model == nil {
		return nil
	}
	return &transaction.Transaction{
		AccountID:       model.AccountID,
		TransactionID:   model.TransactionID,
		OperationTypeID: model.OperationTypeID,
		Amount:          model.Amount,
		EventDate:       model.EventDate,
	}
}
