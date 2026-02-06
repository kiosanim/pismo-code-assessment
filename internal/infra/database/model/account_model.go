package model

type AccountModel struct {
	AccountID      int64  `bun:"account_id,pk,autoincrement"` // Unique identifier of an Account
	DocumentNumber string `bun:"document_number,notnull"`     // Brazilian CPF or CNPJ
}
