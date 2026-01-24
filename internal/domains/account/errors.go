package account

import "errors"

var (
	AccountServiceInvalidParametersError    = errors.New("account service invalid parameter")
	AccountRepositoryInvalidParametersError = errors.New("account repository invalid parameters")
	AccountNotFoundError                    = errors.New("account not found")
)
