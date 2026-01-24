package adapter

import (
	"database/sql"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/config/dto"
	"github.com/uptrace/bun"
)

type ConnectionData struct {
	SqlDB *sql.DB
	BunDB *bun.DB
}

type Connection interface {
	Connect(cfg dto.Configuration) (*sql.DB, error)
}
