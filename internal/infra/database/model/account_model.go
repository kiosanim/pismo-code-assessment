package model

import (
	"github.com/uptrace/bun"
)

type AccountModel struct {
	bun.BaseModel  `bun:"table:accounts,alias:a"`
	AccountID      int64  `bun:"account_id,pk,autoincrement"` // Unique identifier of an Account
	DocumentNumber string `bun:"document_number,notnull"`     // Brazilian CPF or CNPJ
}
