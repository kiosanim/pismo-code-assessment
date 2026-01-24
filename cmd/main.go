package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"log"
	"os"
)

// User model
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            int64  `bun:",pk,autoincrement"`
	Name          string `bun:",notnull"`
	Email         string `bun:",unique"`
}

func main() {
	path, _ := os.Getwd()
	cfg, err := config.LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	// Using pgdriver (recommended)
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(cfg.Database.Dsn),
	))
	if err := sqldb.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer sqldb.Close()
	// Create Bun database instance
	db := bun.NewDB(sqldb, pgdialect.New())
	// Create table
	_, err = db.NewCreateTable().Model((*User)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		panic(err)
	}

	// Insert user
	user := &User{Name: "John Doe", Email: "john@example.com"}
	_, err = db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		panic(err)
	}

	// Select user
	var selectedUser User
	err = db.NewSelect().Model(&selectedUser).Where("email = ?", "john@example.com").Scan(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("User: %+v\n", selectedUser)
}
