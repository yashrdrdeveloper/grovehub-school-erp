package constants

import (
	coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
)

var (
	DependenciesRequired = coreerrors.NewDefinition("AUDIT_DEPENDENCIES_REQUIRED", "audit requires PostgreSQL, Infrastructure and access policy")
	DatabaseRequired = coreerrors.NewDefinition("AUDIT_DATABASE_REQUIRED", "audit migration requires a database")
	ActivityMigrationRequired = coreerrors.NewDefinition("AUDIT_ACTIVITY_MIGRATION_REQUIRED", "audit migrations required")
	RejectionMigrationRequired = coreerrors.NewDefinition("AUDIT_REJECTION_MIGRATION_REQUIRED", "audit rejection migrations required")
	ReceiptMigrationRequired = coreerrors.NewDefinition("AUDIT_RECEIPT_MIGRATION_REQUIRED", "audit receipt migrations required")
	OperatorMigrationRequired = coreerrors.NewDefinition("AUDIT_OPERATOR_MIGRATION_REQUIRED", "operator audit migration required")
	SchoolReceiptMigrationRequired = coreerrors.NewDefinition("AUDIT_SCHOOL_RECEIPT_MIGRATION_REQUIRED", "school receipt migration required")
	MonitoringFailed = coreerrors.NewDefinition("AUDIT_MONITORING_FAILED", "audit monitoring initialization failed")
)
