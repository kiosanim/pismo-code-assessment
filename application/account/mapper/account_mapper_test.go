package mapper

import (
	"testing"

	"github.com/kiosanim/pismo-code-assessment/application/account/dto"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"github.com/stretchr/testify/assert"
)

func TestCreateDTOToEntity(t *testing.T) {
	tests := []struct {
		name     string
		input    dto.CreateAccountRequest
		expected *account.Account
	}{
		{
			name: "should map valid CPF document number",
			input: dto.CreateAccountRequest{
				DocumentNumber: "43904922092",
			},
			expected: &account.Account{
				DocumentNumber: "43904922092",
			},
		},
		{
			name: "should map CNPJ with separators",
			input: dto.CreateAccountRequest{
				DocumentNumber: "12.345.678/0001-90",
			},
			expected: &account.Account{
				DocumentNumber: "12.345.678/0001-90",
			},
		},
		{
			name: "should map empty document number",
			input: dto.CreateAccountRequest{
				DocumentNumber: "",
			},
			expected: &account.Account{
				DocumentNumber: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateDTOToEntity(tt.input)

			assert.NotNil(t, result, "result should not be nil")
			assert.Equal(t, tt.expected.DocumentNumber, result.DocumentNumber, "document number should match")
			assert.Equal(t, int64(0), result.AccountID, "account ID should be zero for new entity")
		})
	}
}

func TestCreateEntityToResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    *account.Account
		expected *dto.CreateAccountResponse
	}{
		{
			name: "should map entity to response with valid data",
			input: &account.Account{
				AccountID:      123,
				DocumentNumber: "43904922092",
			},
			expected: &dto.CreateAccountResponse{
				AccountID:      123,
				DocumentNumber: "43904922092",
			},
		},
		{
			name: "should map entity with zero account ID",
			input: &account.Account{
				AccountID:      0,
				DocumentNumber: "43904922092",
			},
			expected: &dto.CreateAccountResponse{
				AccountID:      0,
				DocumentNumber: "43904922092",
			},
		},
		{
			name: "should map entity with empty document number",
			input: &account.Account{
				AccountID:      456,
				DocumentNumber: "",
			},
			expected: &dto.CreateAccountResponse{
				AccountID:      456,
				DocumentNumber: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateEntityToResponse(tt.input)

			assert.NotNil(t, result, "result should not be nil")
			assert.Equal(t, tt.expected.AccountID, result.AccountID, "account ID should match")
			assert.Equal(t, tt.expected.DocumentNumber, result.DocumentNumber, "document number should match")
		})
	}
}

func TestFindEntityToResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    *account.Account
		expected *dto.FindAccountByIdResponse
	}{
		{
			name: "should map entity to find response with valid data",
			input: &account.Account{
				AccountID:      789,
				DocumentNumber: "43904922092",
			},
			expected: &dto.FindAccountByIdResponse{
				AccountID:      789,
				DocumentNumber: "43904922092",
			},
		},
		{
			name: "should map entity with CNPJ",
			input: &account.Account{
				AccountID:      999,
				DocumentNumber: "42348063000183",
			},
			expected: &dto.FindAccountByIdResponse{
				AccountID:      999,
				DocumentNumber: "42348063000183",
			},
		},
		{
			name: "should map entity with formatted document",
			input: &account.Account{
				AccountID:      111,
				DocumentNumber: "699.415.820-92",
			},
			expected: &dto.FindAccountByIdResponse{
				AccountID:      111,
				DocumentNumber: "699.415.820-92",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindEntityToResponse(tt.input)

			assert.NotNil(t, result, "result should not be nil")
			assert.Equal(t, tt.expected.AccountID, result.AccountID, "account ID should match")
			assert.Equal(t, tt.expected.DocumentNumber, result.DocumentNumber, "document number should match")
		})
	}
}

func TestCreateDTOToEntity_Nil(t *testing.T) {
	// Test that mapper handles empty DTO gracefully
	input := dto.CreateAccountRequest{}
	result := CreateDTOToEntity(input)

	assert.NotNil(t, result, "result should not be nil even with empty DTO")
	assert.Equal(t, "", result.DocumentNumber, "document number should be empty")
	assert.Equal(t, int64(0), result.AccountID, "account ID should be zero")
}

func TestCreateEntityToResponse_WithLargeAccountID(t *testing.T) {
	// Test with large account ID
	input := &account.Account{
		AccountID:      9223372036854775807, // max int64
		DocumentNumber: "69941582092",
	}

	result := CreateEntityToResponse(input)

	assert.NotNil(t, result, "result should not be nil")
	assert.Equal(t, int64(9223372036854775807), result.AccountID, "should handle large account IDs")
}

func TestFindEntityToResponse_WithLargeAccountID(t *testing.T) {
	// Test with large account ID
	input := &account.Account{
		AccountID:      9223372036854775807, // max int64
		DocumentNumber: "69941582092",
	}

	result := FindEntityToResponse(input)

	assert.NotNil(t, result, "result should not be nil")
	assert.Equal(t, int64(9223372036854775807), result.AccountID, "should handle large account IDs")
}
