package model

type OperationTypeModel struct {
	OperationTypeID int64  `bun:"operation_type_id,pk,autoincrement"` // Unique identifier of an OperationType
	Description     string `bun:"description,notnull"`
}
