package transaction

import (
	"time"
)

const (
	Purchase int = iota + 1
	InstallmentPurchase
	Withdrawal
	Payment
)

func IsAValidTransactionType(transactionType int) bool {
	if transactionType >= Purchase && transactionType <= Payment {
		return true
	}
	return false
}

// Transaction represent a transaction
type Transaction struct {
	TransactionID   int64 // Unique identifier of a Transaction
	AccountID       int64
	OperationTypeID int
	Amount          float64
	EventDate       time.Time
}
