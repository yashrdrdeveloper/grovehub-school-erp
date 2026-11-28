// Package mail exposes Integration-owned delivery reconciliation.
package mail

import (
	"context"

	"grovehub/services/integration/client/mail/constants"
	dto "grovehub/services/integration/internal/app/dto/mail"
)

const (
	ReconciliationPath = constants.ReconciliationPath
	ReconciliationScope = constants.ReconciliationScope
	ReconciliationSubject = constants.ReconciliationSubject
)

var ErrInvalidScope = constants.ErrInvalidScope

type (
	SearchRequest = dto.SearchRequest
	Row = dto.Row
	Report = dto.Report
)

// Client is the typed read-only diagnostic boundary. Local callers must be
// trusted operators; remote callers must use authenticated workload HTTP.
// @Client mailClient
type Client interface {
	// @Proxy service.Search
	Search(ctx context.Context, request dto.SearchRequest) (dto.Report, error)
}
