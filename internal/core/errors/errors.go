package errors

import "errors"

var (
	InvalidParametersError                     = errors.New("invalid parameters")
	AccountNotFoundError                       = errors.New("account not found")
	AccountAlreadyExistsForDocumentNumberError = errors.New("an account already exists for this document number")
	ConfigFileNotFountError                    = errors.New("config file not found")
	ConfigFileUnmarshalError                   = errors.New("config unmarshal error")
	TransactionInvalidOperationTypeError       = errors.New("invalid operation type")
	TransactionInvalidAmountNegativeError      = errors.New("invalid amount. must be a positive value")
	TransactionInvalidAccountIDError           = errors.New("invalid account ID")
)
