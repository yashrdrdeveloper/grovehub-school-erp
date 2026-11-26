package mailreconcile

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Row is an operator-only projection. Never select message bodies, addresses,
// token hashes, payload digests or provider errors.
type Row struct {
	DeliveryID string `json:"deliveryId"`
	AccountID *string `json:"accountId,omitempty"`
	Purpose *string `json:"purpose,omitempty"`
	IdentityState string `json:"identityState"`
	IntegrationState string `json:"integrationState"`
	SuppressionReason *string `json:"suppressionReason,omitempty"`
	ReservedAt *time.Time `json:"reservedAt,omitempty"`
	AcceptedAt *time.Time `json:"acceptedAt,omitempty"`
	Disposition string `json:"disposition"`
}

type Report struct {
	TenantID string `json:"tenantId"`
	Rows []Row `json:"rows"`
	NextCursor string `json:"nextCursor,omitempty"`
}

// Read is a cross-schema, read-only operator diagnostic for the shared
// database deployment. It is not a runtime service or an authorization API.
func Read(ctx context.Context, db *gorm.DB, tenant string, after string, limit int) (Report, error) {
	if db == nil || tenant == "" || strings.TrimSpace(tenant) != tenant || limit < 1 || limit > 500 {
		return Report{}, errors.New("invalid reconciliation scope or limit")
	}
	if after == "" {
		after = uuid.Nil.String()
	}
	id, err := uuid.Parse(after)
	if err != nil || id.String() != after {
		return Report{}, errors.New("invalid delivery cursor")
	}
	result := Report{TenantID: tenant, Rows: make([]Row, 0)}
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Raw(reconciliationQuery, tenant, after, tenant, after, limit+1).Scan(&result.Rows).Error
	}, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return Report{}, err
	}
	if len(result.Rows) > limit {
		result.Rows = result.Rows[:limit]
		result.NextCursor = result.Rows[limit-1].DeliveryID
	}
	return result, nil
}
