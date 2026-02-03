package mapper

import (
	"testing"
	"time"

	"github.com/kiosanim/pismo-code-assessment/application/transaction/dto"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	"github.com/stretchr/testify/assert"
)

func TestCreateDTOToEntity(t *testing.T) {
	tests := []struct {
		name     string
		input    dto.CreateTransactionRequest
		expected *transaction.Transaction
	}{
		{
			name: "should map purchase transaction",
			input: dto.CreateTransactionRequest{
				AccountID:       1,
				OperationTypeID: transaction.Purchase,
				Amount:          123.45,
			},
			expected: &transaction.Transaction{
				AccountID:       1,
				OperationTypeID: transaction.Purchase,
				Amount:          123.45,
			},
		},
		{
			name: "should map installment purchase transaction",
			input: dto.CreateTransactionRequest{
				AccountID:       2,
				OperationTypeID: transaction.InstallmentPurchase,
				Amount:          500.00,
			},
			expected: &transaction.Transaction{
				AccountID:       2,
				OperationTypeID: transaction.InstallmentPurchase,
				Amount:          500.00,
			},
		},
		{
			name: "should map withdrawal transaction",
			input: dto.CreateTransactionRequest{
				AccountID:       3,
				OperationTypeID: transaction.Withdrawal,
				Amount:          75.50,
			},
			expected: &transaction.Transaction{
				AccountID:       3,
				OperationTypeID: transaction.Withdrawal,
				Amount:          75.50,
			},
		},
		{
			name: "should map payment transaction",
			input: dto.CreateTransactionRequest{
				AccountID:       4,
				OperationTypeID: transaction.Payment,
				Amount:          1000.00,
			},
			expected: &transaction.Transaction{
				AccountID:       4,
				OperationTypeID: transaction.Payment,
				Amount:          1000.00,
			},
		},
		{
			name: "should map transaction with zero values",
			input: dto.CreateTransactionRequest{
				AccountID:       0,
				OperationTypeID: 0,
				Amount:          0,
			},
			expected: &transaction.Transaction{
				AccountID:       0,
				OperationTypeID: 0,
				Amount:          0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateDTOToEntity(tt.input)

			assert.NotNil(t, result, "result should not be nil")
			assert.Equal(t, tt.expected.AccountID, result.AccountID, "account ID should match")
			assert.Equal(t, tt.expected.OperationTypeID, result.OperationTypeID, "operation type ID should match")
			assert.Equal(t, tt.expected.Amount, result.Amount, "amount should match")
			assert.Equal(t, int64(0), result.TransactionID, "transaction ID should be zero for new entity")
			assert.True(t, result.EventDate.IsZero(), "event date should be zero time")
		})
	}
}

func TestEntityToResponse(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    *transaction.Transaction
		expected *dto.CreateTransactionResponse
	}{
		{
			name: "should map entity to response for purchase",
			input: &transaction.Transaction{
				TransactionID:   100,
				AccountID:       1,
				OperationTypeID: transaction.Purchase,
				Amount:          -123.45, // Negative in storage
				EventDate:       now,
			},
			expected: &dto.CreateTransactionResponse{
				TransactionID:   100,
				AccountID:       1,
				OperationTypeID: transaction.Purchase,
				Amount:          -123.45,
			},
		},
		{
			name: "should map entity to response for payment",
			input: &transaction.Transaction{
				TransactionID:   200,
				AccountID:       2,
				OperationTypeID: transaction.Payment,
				Amount:          500.00, // Positive in storage
				EventDate:       now,
			},
			expected: &dto.CreateTransactionResponse{
				TransactionID:   200,
				AccountID:       2,
				OperationTypeID: transaction.Payment,
				Amount:          500.00,
			},
		},
		{
			name: "should map entity with zero transaction ID",
			input: &transaction.Transaction{
				TransactionID:   0,
				AccountID:       5,
				OperationTypeID: transaction.Withdrawal,
				Amount:          -50.00,
				EventDate:       now,
			},
			expected: &dto.CreateTransactionResponse{
				TransactionID:   0,
				AccountID:       5,
				OperationTypeID: transaction.Withdrawal,
				Amount:          -50.00,
			},
		},
		{
			name: "should map entity with large values",
			input: &transaction.Transaction{
				TransactionID:   9999999,
				AccountID:       9999,
				OperationTypeID: transaction.Payment,
				Amount:          999999.99,
				EventDate:       now,
			},
			expected: &dto.CreateTransactionResponse{
				TransactionID:   9999999,
				AccountID:       9999,
				OperationTypeID: transaction.Payment,
				Amount:          999999.99,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EntityToResponse(tt.input)

			assert.NotNil(t, result, "result should not be nil")
			assert.Equal(t, tt.expected.TransactionID, result.TransactionID, "transaction ID should match")
			assert.Equal(t, tt.expected.AccountID, result.AccountID, "account ID should match")
			assert.Equal(t, tt.expected.OperationTypeID, result.OperationTypeID, "operation type ID should match")
			assert.Equal(t, tt.expected.Amount, result.Amount, "amount should match")
		})
	}
}

func TestCreateDTOToEntity_EmptyDTO(t *testing.T) {
	input := dto.CreateTransactionRequest{}
	result := CreateDTOToEntity(input)

	assert.NotNil(t, result, "result should not be nil even with empty DTO")
	assert.Equal(t, int64(0), result.AccountID, "account ID should be zero")
	assert.Equal(t, 0, result.OperationTypeID, "operation type ID should be zero")
	assert.Equal(t, 0.0, result.Amount, "amount should be zero")
	assert.Equal(t, int64(0), result.TransactionID, "transaction ID should be zero")
}

func TestEntityToResponse_WithLargeAccountID(t *testing.T) {
	input := &transaction.Transaction{
		TransactionID:   9223372036854775807, // max int64
		AccountID:       9223372036854775807,
		OperationTypeID: transaction.Purchase,
		Amount:          123456789.99,
		EventDate:       time.Now(),
	}

	result := EntityToResponse(input)

	assert.NotNil(t, result, "result should not be nil")
	assert.Equal(t, int64(9223372036854775807), result.TransactionID, "should handle large transaction IDs")
	assert.Equal(t, int64(9223372036854775807), result.AccountID, "should handle large account IDs")
}

func TestCreateDTOToEntity_AllOperationTypes(t *testing.T) {
	operationTypes := []struct {
		name string
		code int
	}{
		{"Purchase", transaction.Purchase},
		{"InstallmentPurchase", transaction.InstallmentPurchase},
		{"Withdrawal", transaction.Withdrawal},
		{"Payment", transaction.Payment},
	}

	for _, op := range operationTypes {
		t.Run(op.name, func(t *testing.T) {
			input := dto.CreateTransactionRequest{
				AccountID:       1,
				OperationTypeID: op.code,
				Amount:          100.00,
			}

			result := CreateDTOToEntity(input)

			assert.NotNil(t, result, "result should not be nil")
			assert.Equal(t, op.code, result.OperationTypeID, "operation type should match")
		})
	}
}

func TestEntityToResponse_NegativeAndPositiveAmounts(t *testing.T) {
	tests := []struct {
		name          string
		amount        float64
		expectedValue float64
	}{
		{
			name:          "positive amount",
			amount:        100.00,
			expectedValue: 100.00,
		},
		{
			name:          "negative amount",
			amount:        -100.00,
			expectedValue: -100.00,
		},
		{
			name:          "zero amount",
			amount:        0.00,
			expectedValue: 0.00,
		},
		{
			name:          "decimal amount",
			amount:        123.456789,
			expectedValue: 123.456789,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &transaction.Transaction{
				TransactionID:   1,
				AccountID:       1,
				OperationTypeID: transaction.Purchase,
				Amount:          tt.amount,
				EventDate:       time.Now(),
			}

			result := EntityToResponse(input)

			assert.NotNil(t, result, "result should not be nil")
			assert.Equal(t, tt.expectedValue, result.Amount, "amount should be preserved")
		})
	}
}