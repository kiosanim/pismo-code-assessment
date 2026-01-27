package transaction

import (
	"testing"
	"time"
)

func TestIsAValidTransactionType(t *testing.T) {
	type args struct {
		transaction   *Transaction
		operationType int
	}
	tests := []struct {
		name          string
		args          args
		shouldBeValid bool
	}{
		{name: "must have a valid transaction type", args: args{transaction: &Transaction{TransactionID: 1, AccountID: 1, OperationTypeID: Purchase, Amount: 10.0, EventDate: time.Now()}, operationType: Purchase}, shouldBeValid: true},
		{name: "must have an invalid transaction type", args: args{transaction: &Transaction{TransactionID: 1, AccountID: 1, OperationTypeID: 0, Amount: 10.0, EventDate: time.Now()}, operationType: Purchase}, shouldBeValid: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAValidOperationType(tt.args.transaction.OperationTypeID)
			if got != tt.shouldBeValid {
				t.Errorf("IsAValidOperationType() = %v, want %v", got, tt.shouldBeValid)
			}
		})
	}
}
