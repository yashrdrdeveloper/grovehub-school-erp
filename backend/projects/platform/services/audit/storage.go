package audit

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// StorageTable contains database-wide operational estimates, not tenant data.
// PostgreSQL statistics may be stale or reset; they are not exact row counts.
type StorageTable struct {
	Name          string `json:"name"`
	TableBytes    int64  `json:"tableBytes"`
	IndexBytes    int64  `json:"indexBytes"`
	TotalBytes    int64  `json:"totalBytes"`
	EstimatedRows int64  `json:"estimatedRows"`
	DeadRows      int64  `json:"estimatedDeadRows"`
}

// InspectStorage is an operator-only, read-only catalog inspection.
// It does not scan activity payloads or modify history.
func InspectStorage(ctx context.Context, db *gorm.DB, expectedDatabase string) ([]StorageTable, error) {
	if db == nil || strings.TrimSpace(expectedDatabase) == "" || strings.TrimSpace(expectedDatabase) != expectedDatabase {
		return nil, errors.New("audit storage inspection requires an exact database name")
	}
	pool, err := db.DB()
	if err != nil {
		return nil, err
	}
	tx, err := pool.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var database string
	if err := tx.QueryRowContext(ctx, "SELECT current_database()").Scan(&database); err != nil {
		return nil, err
	}
	if database != expectedDatabase {
		return nil, errors.New("configured Audit database does not match expected database")
	}
	rows, err := tx.QueryContext(ctx, `
		SELECT relname, pg_table_size(relid), pg_indexes_size(relid),
		       pg_total_relation_size(relid), n_live_tup, n_dead_tup
		FROM pg_stat_user_tables
		WHERE schemaname = 'audit'
		  AND relname IN ('activities', 'event_receipts', 'rejected_events', 'replay_runs')
		ORDER BY relname`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]StorageTable, 0, 4)
	for rows.Next() {
		var table StorageTable
		if err := rows.Scan(&table.Name, &table.TableBytes, &table.IndexBytes, &table.TotalBytes, &table.EstimatedRows, &table.DeadRows); err != nil {
			return nil, err
		}
		result = append(result, table)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(result) != 4 {
		return nil, errors.New("Audit storage schema is incomplete")
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}
