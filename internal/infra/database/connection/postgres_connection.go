package connection

import (
	"database/sql"
	"github.com/kiosanim/pismo-code-assessment/internal/core/adapter"
	"github.com/kiosanim/pismo-code-assessment/internal/core/config"
	"github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	_ "github.com/lib/pq"
)

type PostgresConnection struct {
	configuration  *config.Configuration
	connectionData *adapter.DatabaseConnectionData
}

func NewPostgresConnection(configuration *config.Configuration) *PostgresConnection {
	return &PostgresConnection{
		configuration: configuration,
	}
}

func (p *PostgresConnection) Connect() (*adapter.DatabaseConnectionData, error) {
	db, err := sql.Open("postgres", p.configuration.Database.URL)
	if err != nil {
		return nil, errors.DatabaseConnectionFailedError
	}
	err = db.Ping()
	if err != nil {
		return nil, errors.DatabaseConnectionValidationFailedError
	}
	p.connectionData = &adapter.DatabaseConnectionData{db}
	return p.connectionData, nil
}
