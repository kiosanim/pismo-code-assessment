package model

import (
	"github.com/uptrace/bun"
	"time"
)

type TransactionModel struct {
	bun.BaseModel   `bun:"table:transactions,alias:t"`
	TransactionID   int64     `bun:"transaction_id,pk,autoincrement"` // Unique identifier of an Account
	AccountID       int64     `bun:"account_id,notnull"`
	OperationTypeID int       `bun:"operation_type_id,notnull"`
	Amount          float64   `bun:"amount,notnull"`
	EventDate       time.Time `bun:"event_date,notnull"`

	Account       *AccountModel       `bun:"rel:belongs-to,join:account_id=account_id"`
	OperationType *OperationTypeModel `bun:"rel:belongs-to,join:operation_type_id=operation_type_id"`
}
