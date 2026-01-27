package account

import (
	"github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/paemuri/brdoc"
	"regexp"
)

var (
	documentNumberRegex = regexp.MustCompile(`\D`)
)

// Account represent a customer account
type Account struct {
	AccountID      int64  // Unique identifier of an Account
	DocumentNumber string // Brazilian CPF or CNPJ
}

// IsValidDocumentNumber validate if a user DocumentNumber is a Brazilian CPF or CNPJ
func IsValidDocumentNumber(documentNumber string) error {
	if brdoc.IsCPF(documentNumber) || brdoc.IsCNPJ(documentNumber) {
		return nil
	}
	return errors.InvalidParametersError
}

// SanitizeDocumentNumber removes all characters except digits from a DocumentNumber
func SanitizeDocumentNumber(documentNumber string) string {
	return documentNumberRegex.ReplaceAllString(documentNumber, "")
}
