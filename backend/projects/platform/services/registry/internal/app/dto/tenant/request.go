package tenant

import "grovehub/services/registry/internal/app/constants/status"

type CreateInput struct {
	ExternalReference string `json:"externalReference"`
}

type StatusInput struct {
	ExpectedVersion int64         `json:"expectedVersion"`
	Status          status.Status `json:"status"`
}
