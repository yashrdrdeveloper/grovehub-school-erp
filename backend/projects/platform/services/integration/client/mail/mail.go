// Package mail exposes Integration-owned delivery reconciliation.
package mail

import (
	"context"
	"errors"
	"time"
)

const ReconciliationPath = "/integrations/restricted/mail/reconciliation/search"
const ReconciliationScope = "integration.mail.reconciliation.read"
const ReconciliationSubject = "platform-mail-reconciliation-reader"

var ErrInvalidScope = errors.New("invalid reconciliation scope or cursor")

type SearchRequest struct {
	TenantID string `json:"tenantId"`
	After string `json:"after,omitempty"`
	Limit int `json:"limit"`
}

// Row contains only operator diagnostic fields. ACCEPTED is provider
// acceptance, not inbox delivery.
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

// Client is the typed read-only diagnostic boundary. Local callers must be
// trusted operators; remote callers must use authenticated workload HTTP.
type Client interface {
	Search(context.Context, SearchRequest) (Report, error)
}
