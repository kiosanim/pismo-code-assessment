package dto

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" validate:"required"`
}

type CreateAccountResponse struct {
	AccountID      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

type FindAccountByIdRequest struct {
	AccountID int64 `uri:"account_id" binding:"required" validate:"required,number,gt=0"`
}

type FindAccountByIdResponse struct {
	AccountID      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}
