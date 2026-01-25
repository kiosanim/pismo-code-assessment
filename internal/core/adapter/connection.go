package adapter

import (
	"context"
	"database/sql"
	"errors"
	"github.com/uptrace/bun"
)

var (
	ConnectionFailedError = errors.New("Failed to connect to database")
)

type ConnectionData struct {
	SqlDB *sql.DB
	BunDB *bun.DB
}

type Connection interface {
	Connect(ctx context.Context) (*ConnectionData, error)
}
