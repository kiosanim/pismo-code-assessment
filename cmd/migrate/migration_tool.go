package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"github.com/kiosanim/pismo-code-assessment/internal/infra/config"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

// func main() {
// path, _ := os.Getwd()
// cfg, err := config.LoadConfig(path)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
// db, err := sql.Open("postgres", cfg.Database.Dsn)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
// goose.SetDialect("postgres")
// fmt.Println("Connected to database")
// migrationsFolder := "internal/infra/database/migrations"
// fmt.Println("Migrations folder is", migrationsFolder)
// fmt.Println("Running migrations")
// //err = goose.UpContext(context.Background(), db, migrationsFolder)
// //if err != nil {
// //	log.Fatalf("failed to run migrations: %v\n", err)
// //}
// //
// err = goose.DownToContext(context.Background(), db, migrationsFolder, 0)
//
//	if err != nil {
//		log.Fatalf("failed to run migrations: %v\n", err)
//	}
//
// fmt.Println("Migrations applied successfully")
// }

const MigrationsFolder = "internal/infra/database/migrations"

func up(ctx context.Context, db *sql.DB) *cobra.Command {
	upCmd := &cobra.Command{
		Use:   "up",
		Short: "Run up migration",
		Run: func(cmd *cobra.Command, args []string) {
			err := goose.UpContext(ctx, db, MigrationsFolder)
			if err != nil {
				log.Fatalf("failed to run migrations: %v\n", err)
			}
		},
	}
	return upCmd
}

func destroy(ctx context.Context, db *sql.DB) *cobra.Command {
	downCmd := &cobra.Command{
		Use:   "down",
		Short: "Rollback All Migrations",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("⚠️ Are you sure you want to destroy every table? (y/N): ")
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" && answer != "yes" {
				fmt.Println("❌ Aborted.")
				return
			}
			err := goose.DownToContext(ctx, db, MigrationsFolder, 0)
			if err != nil {
				log.Fatalf("failed to run migrations: %v\n", err)
			}
		},
	}
	return downCmd
}

func status(ctx context.Context, db *sql.DB) *cobra.Command {
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Status Migrations",
		Run: func(cmd *cobra.Command, args []string) {
			err := goose.StatusContext(ctx, db, MigrationsFolder)
			if err != nil {
				log.Fatalf("failed to run migrations: %v\n", err)
			}
		},
	}
	return statusCmd
}

func main() {
	path, _ := os.Getwd()
	cfg, err := config.LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	db, err := sql.Open("postgres", cfg.Database.Dsn)

	if err != nil {
		log.Fatal(err)
	}
	goose.SetDialect("postgres")
	fmt.Println("Connected to database")
	migrationsFolder := "internal/infra/database/migrations"
	fmt.Println("Migrations folder is", migrationsFolder)
	var rootCmd = &cobra.Command{
		Use:   "migration-tool",
		Short: "Simple CLI to up/down migrations",
	}
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(up(ctx, db), destroy(ctx, db), status(ctx, db))
	err = rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
