package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	"grovehub/operations/mailreconcile"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "identity-mail-reconcile:", err)
		os.Exit(1)
	}
}

func run() error {
	tenant := flag.String("tenant-id", "", "Exact tenant identifier")
	database := flag.String("database", "", "Expected database name")
	after := flag.String("after", "", "Continuation delivery UUID")
	limit := flag.Int("limit", 100, "Maximum rows (1-500)")
	flag.Parse()
	if flag.NArg() != 0 || *tenant == "" || strings.TrimSpace(*tenant) != *tenant || *database == "" {
		return errors.New("exact --tenant-id and --database are required")
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return errors.New("DATABASE_URL is required")
	}
	db, err := coredatabase.NewDB(dsn)
	if err != nil {
		return errors.New("could not connect to configured database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var current string
	if err := db.WithContext(ctx).Raw("SELECT current_database()").Scan(&current).Error; err != nil {
		return errors.New("could not verify database name")
	}
	if current != *database {
		return errors.New("configured database does not match --database")
	}
	report, err := mailreconcile.Read(ctx, db, *tenant, *after, *limit)
	if err != nil {
		return errors.New("reconciliation query failed; verify scope and migrations")
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}
