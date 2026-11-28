package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	bootstrap "grovehub/launcher"
	"grovehub/services/audit"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	database := flag.String("database", "", "Exact expected Audit database name")
	flag.Parse()
	if flag.NArg() != 0 || *database == "" {
		return errors.New("database is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	db, closeDB, err := bootstrap.NewAuditPostgresDB(ctx, nil)
	if err != nil {
		return errors.New("could not connect to configured Audit database")
	}
	defer closeDB()
	tables, err := audit.InspectStorage(ctx, db, *database)
	if err != nil {
		return errors.New("Audit storage inspection failed; verify database name, schema and permissions")
	}
	return json.NewEncoder(os.Stdout).Encode(tables)
}
