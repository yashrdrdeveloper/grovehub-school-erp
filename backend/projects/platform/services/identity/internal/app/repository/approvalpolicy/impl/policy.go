package impl

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	authorityconstants "grovehub/services/identity/internal/app/repository/authority/constants"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/repository/approvalpolicy"
	"grovehub/services/identity/internal/app/repository/approvalpolicy/constants"
	"grovehub/services/identity/internal/app/repository/approvalpolicy/helpers"
	"grovehub/services/identity/internal/app/repository/authority"
	"grovehub/services/identity/internal/app/repository/rbacmigration"
)

type Provisioner struct {
	db *gorm.DB
	readiness rbacmigration.Readiness
	administrators authority.AdministratorReader
}

var _ approvalpolicy.Provisioner = (*Provisioner)(nil)

func New(db *gorm.DB, readiness rbacmigration.Readiness, administrators authority.AdministratorReader) approvalpolicy.Provisioner {
	return &Provisioner{db: db, readiness: readiness, administrators: administrators}
}

func (writer *Provisioner) Apply(ctx context.Context, manifest approvalpolicy.Manifest, operatorID string, digest string) (approvalpolicy.Receipt, error) {
	if err := helpers.Validate(manifest, operatorID, digest); err != nil {
		return approvalpolicy.Receipt{}, err
	}
	if writer.db == nil || writer.readiness == nil || writer.administrators == nil {
		return approvalpolicy.Receipt{}, coredatabase.ErrMissingDatabase
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return approvalpolicy.Receipt{}, err
	}
	receipt := approvalpolicy.Receipt{ReceiptID: id.String(), Policies: make([]approvalpolicy.Activated, 0, len(manifest.Policies))}
	err = coredatabase.RunTransaction(ctx, writer.db, func(tx context.Context) error {
		db := coredatabase.TransactionDB(tx, writer.db)
		// Offline only; caller must not already hold tenant authority locks.
		if err := db.Exec(authorityconstants.LockGlobalExclusive).Error; err != nil {
			return err
		}
		if err := db.Exec(constants.Lock).Error; err != nil {
			return err
		}
		if err := writer.readiness.Require(tx, manifest.TenantID); err != nil {
			return err
		}
		count, err := writer.administrators.Count(tx, manifest.TenantID)
		if err != nil {
			return err
		}
		if count < 1 {
			return authority.ErrBootstrapRequired
		}
		if err := db.Raw(constants.Author, manifest.TenantID, manifest.AuthorID).Scan(&count).Error; err != nil {
			return err
		}
		if count != 1 {
			return approvalpolicy.ErrInvalid
		}
		for _, revision := range manifest.Policies {
			activated, err := activate(db, manifest, revision)
			if err != nil {
				return err
			}
			receipt.Policies = append(receipt.Policies, activated)
		}
		payload, err := json.Marshal(receipt)
		if err != nil {
			return err
		}
		return db.Exec(constants.Receipt, receipt.ReceiptID, manifest.TenantID, operatorID, manifest.AuthorID, digest, string(payload)).Error
	})
	if err != nil {
		return approvalpolicy.Receipt{}, err
	}
	return receipt, nil
}

func activate(db *gorm.DB, manifest approvalpolicy.Manifest, revision approvalpolicy.Revision) (approvalpolicy.Activated, error) {
	var version int64
	if err := db.Raw(constants.Version, manifest.TenantID, revision.Type).Scan(&version).Error; err != nil {
		return approvalpolicy.Activated{}, err
	}
	if version != revision.ExpectedVersion {
		return approvalpolicy.Activated{}, approvalpolicy.ErrConflict
	}
	for _, step := range revision.Steps {
		definition, known := permission.Find(step.Permission)
		if !known {
			return approvalpolicy.Activated{}, approvalpolicy.ErrInvalid
		}
		var count int64
		if err := db.Raw(constants.Eligible, manifest.TenantID, step.Permission, definition.Feature, definition.Feature, step.RoleID, step.RoleID).Scan(&count).Error; err != nil {
			return approvalpolicy.Activated{}, err
		}
		// Two accounts allow requester exclusion. Assignment requests can
		// still exclude a second beneficiary; this is not a liveness guarantee.
		if count < 2 {
			return approvalpolicy.Activated{}, approvalpolicy.ErrReviewers
		}
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return approvalpolicy.Activated{}, err
	}
	if err := db.Exec(constants.Deactivate, manifest.TenantID, revision.Type).Error; err != nil {
		return approvalpolicy.Activated{}, err
	}
	next := approvalpolicy.Activated{Type: revision.Type, PolicyID: id.String(), Version: version+1}
	if err := db.Exec(constants.Insert, next.PolicyID, manifest.TenantID, next.Type, next.Version, manifest.AuthorID).Error; err != nil {
		return approvalpolicy.Activated{}, err
	}
	for index, step := range revision.Steps {
		if err := db.Exec(constants.InsertStep, manifest.TenantID, next.PolicyID, index+1, step.Permission, step.RoleID).Error; err != nil {
			return approvalpolicy.Activated{}, err
		}
	}
	return next, nil
}
