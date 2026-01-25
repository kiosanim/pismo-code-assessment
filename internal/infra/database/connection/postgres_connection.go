package connection

import (
	"context"
	"database/sql"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	"github.com/kiosanim/pismo-code-assessment/internal/core/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type PostgresConnection struct {
	configuration  *config.Configuration
	connectionData *adapter.ConnectionData
}

func NewPostgresConnection(configuration *config.Configuration) *PostgresConnection {
	return &PostgresConnection{
		configuration: configuration,
	}
}

func (p *PostgresConnection) Connect(ctx context.Context) (*adapter.ConnectionData, error) {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(p.configuration.Database.Dsn),
	))
	if err := sqlDb.Ping(); err != nil {
		return nil, adapter.ConnectionFailedError
	}
	bunDb := bun.NewDB(sqlDb, pgdialect.New())
	p.connectionData = &adapter.ConnectionData{
		SqlDB: sqlDb,
		BunDB: bunDb,
	}
	return p.connectionData, nil
}
