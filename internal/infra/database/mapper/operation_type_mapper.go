package mapper

import (
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/database/model"
)

func ToOperationTypeEntity(model *model.OperationTypeModel) *transaction.OperationType {
	if model == nil {
		return nil
	}
	return &transaction.OperationType{
		OperationTypeID: model.OperationTypeID,
		Description:     model.Description,
	}
}
