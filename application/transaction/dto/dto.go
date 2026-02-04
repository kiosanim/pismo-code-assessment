package dto

type TransactionDTO struct {
	TransactionID   int64   `json:"transaction_id"`
	AccountID       int64   `json:"account_id"`
	OperationTypeID int     `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}
type CreateTransactionRequest struct {
	AccountID       int64   `json:"account_id"`
	OperationTypeID int     `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}

type CreateTransactionResponse struct {
	Transaction TransactionDTO `json:"transaction"`
}

type ListTransactionsRequest struct {
	Page int64 `json:"page" validate:"min=0"`
	Size int64 `json:"size" validate:"min=1"`
}

type ListTransactionsResponse struct {
	Transactions []TransactionDTO `json:"transactions"`
}

type FindTransactionByIdRequest struct {
	TransactionID int64 `uri:"transaction_id" binding:"required,gt=0"`
}

type FindTransactionByIdResponse struct {
	Transaction TransactionDTO `json:"transaction"`
}
