package transaction

import "errors"

var (
	TransactionServiceInvalidParametersError    = errors.New("transaction service invalid parameter")
	TransactionRepositoryInvalidParametersError = errors.New("transaction repository invalid parameters")
)
