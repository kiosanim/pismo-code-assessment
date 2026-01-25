package mapper

import (
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/model"
)

func ToAccountModel(entity *account.Account) *model.AccountModel {
	if entity == nil {
		return nil
	}
	return &model.AccountModel{
		AccountID:      entity.AccountID,
		DocumentNumber: entity.DocumentNumber,
	}
}

func ToAccountEntity(model *model.AccountModel) *account.Account {
	if model == nil {
		return nil
	}
	return &account.Account{
		AccountID:      model.AccountID,
		DocumentNumber: model.DocumentNumber,
	}
}
