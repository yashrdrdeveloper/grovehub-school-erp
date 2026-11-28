// Package rbacmigrationcli is the operator-only entrypoint for staged RBAC
// migration. It is never imported by the running school platform.
package rbacmigrationcli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"grovehub/services/identity/internal/app/repository/approvalpolicy"
	policyimpl "grovehub/services/identity/internal/app/repository/approvalpolicy/impl"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"

	migrationimpl "grovehub/services/identity/internal/app/repository/rbacmigration/impl"
	"grovehub/services/identity/internal/app/services/rbacmigration"
	"grovehub/services/identity/internal/app/services/rbacmigration/operations/plan"
	"grovehub/services/identity/rbacmigrationcli/helpers"
)

type options struct {
	operation string
	tenant string
	database string
	operator string
	input string
	digest string
	confirm bool
}

func Run(args []string, output io.Writer, diagnostics io.Writer) error {
	if len(args) == 0 {
		return errors.New("expected plan, backfill, bootstrap or policies")
	}
	opts, err := parse(args, diagnostics)
	if err != nil {
		return err
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return errors.New("DATABASE_URL is missing")
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	var database string
	if err := db.WithContext(ctx).Raw("SELECT current_database()").Scan(&database).Error; err != nil {
		return err
	}
	if database != opts.database {
		return errors.New("configured database does not match --database")
	}
	if opts.operation == "plan" {
		inventory, err := migrationimpl.New(db).Inventory(ctx, opts.tenant)
		if err != nil {
			return err
		}
		result, err := plan.Build(inventory)
		if err != nil {
			return err
		}
		encoder := json.NewEncoder(output)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}
	return apply(ctx, db, opts, output)
}

func parse(args []string, diagnostics io.Writer) (options, error) {
	opts := options{operation: args[0]}
	if opts.operation != "plan" && opts.operation != "backfill" && opts.operation != "bootstrap" && opts.operation != "policies" {
		return options{}, errors.New("unknown migration operation")
	}
	flags := flag.NewFlagSet(opts.operation, flag.ContinueOnError)
	flags.SetOutput(diagnostics)
	flags.StringVar(&opts.tenant, "tenant-id", "", "Exact tenant identifier")
	flags.StringVar(&opts.database, "database", "", "Expected PostgreSQL database name")
	flags.StringVar(&opts.operator, "operator-id", "", "Operator audit UUID")
	flags.StringVar(&opts.input, "input", "", "Reviewed plan or manifest JSON")
	flags.StringVar(&opts.digest, "sha256", "", "SHA-256 of reviewed input bytes")
	flags.BoolVar(&opts.confirm, "confirm-offline-maintenance", false, "Confirm all application writers stopped and backup verified")
	if err := flags.Parse(args[1:]); err != nil {
		return options{}, err
	}
	if flags.NArg() != 0 || opts.tenant == "" || strings.TrimSpace(opts.tenant) != opts.tenant || opts.database == "" {
		return options{}, errors.New("exact --tenant-id and --database are required")
	}
	if opts.operation == "plan" {
		if opts.confirm || opts.operator != "" || opts.input != "" || opts.digest != "" {
			return options{}, errors.New("plan is read-only and does not accept write flags")
		}
		return opts, nil
	}
	id, err := uuid.Parse(opts.operator)
	if err != nil || id == uuid.Nil || id.String() != opts.operator || !opts.confirm || opts.input == "" || opts.digest == "" {
		return options{}, errors.New("writes require --operator-id, --input, --sha256 and --confirm-offline-maintenance")
	}
	return opts, nil
}

func apply(ctx context.Context, db *gorm.DB, opts options, output io.Writer) error {
	data, err := helpers.ReadReviewed(opts.input, opts.digest)
	if err != nil {
		return err
	}
	if opts.operation == "policies" {
		var manifest approvalpolicy.Manifest
		if err := helpers.Decode(data, &manifest); err != nil {
			return err
		}
		if manifest.TenantID != opts.tenant {
			return errors.New("reviewed input tenant does not match --tenant-id")
		}
		receipt, err := policyimpl.New(db, migrationimpl.NewReadiness(db), authorityimpl.NewAdministratorReader(db)).Apply(ctx, manifest, opts.operator, opts.digest)
		if err != nil {
			return err
		}
		if err := json.NewEncoder(output).Encode(receipt); err != nil {
			return fmt.Errorf("policies committed (receipt %s), but output failed: %w", receipt.ReceiptID, err)
		}
		return nil
	}
	var reviewed rbacmigration.Plan
	var manifest rbacmigration.Manifest
	if opts.operation == "backfill" {
		if err := helpers.Decode(data, &reviewed); err != nil {
			return err
		}
	} else {
		if err := helpers.Decode(data, &manifest); err != nil {
			return err
		}
		reviewed = manifest.Plan
	}
	if reviewed.TenantID != opts.tenant {
		return errors.New("reviewed input tenant does not match --tenant-id")
	}
	receiptID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	var result any
	err = coredatabase.RunTransaction(ctx, db, func(tx context.Context) error {
		if opts.operation == "backfill" {
			result, err = migrationimpl.NewBackfill(db).Apply(tx, reviewed, opts.operator)
		} else {
			result, err = migrationimpl.NewBootstrap(db).Apply(tx, manifest, opts.operator)
		}
		if err != nil {
			return err
		}
		payload, err := json.Marshal(result)
		if err != nil {
			return err
		}
		return coredatabase.TransactionDB(tx, db).Exec(
			"INSERT INTO identity.rbac_migration_receipts (receipt_id, tenant_id, operator_id, operation, input_sha256, result) VALUES (?, ?, ?, ?, ?, CAST(? AS jsonb))",
			receiptID.String(), opts.tenant, opts.operator, strings.ToUpper(opts.operation), opts.digest, string(payload),
		).Error
	})
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(output, "%s committed; receipt_id=%s tenant=%s\n", opts.operation, receiptID.String(), opts.tenant)
	if err != nil {
		return fmt.Errorf("migration committed (receipt %s), but output failed: %w", receiptID.String(), err)
	}
	return nil
}
