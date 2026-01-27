package model

import "github.com/uptrace/bun"

type OperationTypeModel struct {
	bun.BaseModel   `bun:"table:operation_types,alias:t"`
	OperationTypeID int64  `bun:"operation_type_id,pk,autoincrement"` // Unique identifier of an OperationType
	Description     string `bun:"description,notnull"`
}
