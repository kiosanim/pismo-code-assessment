package adapter

import (
	"database/sql"
)

type DatabaseConnectionData struct {
	Db *sql.DB
}

type DatabaseConnection interface {
	Connect() (*DatabaseConnectionData, error)
}
