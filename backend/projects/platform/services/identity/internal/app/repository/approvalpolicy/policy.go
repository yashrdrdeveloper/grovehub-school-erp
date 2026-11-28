package approvalpolicy

import (
	"context"
	"errors"

	kind "grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/permission"
)

var (
	ErrInvalid = errors.New("invalid approval policy manifest")
	ErrConflict = errors.New("approval policy version conflict")
	ErrReviewers = errors.New("insufficient eligible independent reviewers")
)

type Step struct {
	Permission permission.Code `json:"permissionCode"`
	RoleID *string `json:"roleId,omitempty"`
}

type Revision struct {
	Type kind.RequestType `json:"requestType"`
	ExpectedVersion int64 `json:"expectedVersion"`
	Steps []Step `json:"steps"`
}

type Manifest struct {
	TenantID string `json:"tenantId"`
	AuthorID string `json:"authorId"`
	Policies []Revision `json:"policies"`
}

type Activated struct {
	Type kind.RequestType `json:"requestType"`
	PolicyID string `json:"policyId"`
	Version int64 `json:"version"`
}

type Receipt struct {
	ReceiptID string `json:"receiptId"`
	Policies []Activated `json:"policies"`
}

// Provisioner is exclusively an offline operator persistence boundary.
type Provisioner interface {
	Apply(ctx context.Context, manifest Manifest, operatorID string, digest string) (Receipt, error)
}
