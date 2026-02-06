package dto

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" binding:"required"`
}

type CreateAccountResponse struct {
	AccountID      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

type FindAccountByIdRequest struct {
	AccountID int64 `uri:"account_id" binding:"required,gt=0"`
}

type FindAccountByIdResponse struct {
	AccountID      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

type ListAccountsRequest struct {
	Limit  int64 `uri:"limit" binding:"required,gt=0"`
	Cursor int64 `uri:"cursor" binding:"required,gt=0"`
}

type ListAccountsResponse struct {
	Accounts []AccountDTO `json:"accounts"`
	Limit    int64        `uri:"limit" binding:"required,gt=0"`
	Cursor   int64        `uri:"cursor" binding:"required,gt=0"`
}

type AccountDTO struct {
	AccountID      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}
