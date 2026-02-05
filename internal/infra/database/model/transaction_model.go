package model

import (
	"time"
)

type TransactionModel struct {
	TransactionID   int64     `bun:"transaction_id,pk,autoincrement"` // Unique identifier of an Account
	AccountID       int64     `bun:"account_id,notnull"`
	OperationTypeID int       `bun:"operation_type_id,notnull"`
	Amount          float64   `bun:"amount,notnull"`
	EventDate       time.Time `bun:"event_date,notnull"`
}
