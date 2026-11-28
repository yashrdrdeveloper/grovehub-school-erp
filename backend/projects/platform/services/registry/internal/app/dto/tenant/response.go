package tenant

import (
	"time"

	"grovehub/services/registry/internal/app/constants/status"
	"grovehub/services/registry/internal/app/entity"
)

type Response struct {
	ProjectID string `json:"projectId"`
	TenantID string `json:"tenantId"`
	ExternalReference string `json:"externalReference"`
	SessionGeneration int64 `json:"sessionGeneration"`
	Status status.Status `json:"status"`
	Version int64 `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func FromEntity(value entity.Tenant) Response {
	return Response{
		ProjectID: value.ProjectID,
		TenantID: value.TenantID,
		ExternalReference: value.ExternalReference,
		SessionGeneration: value.SessionGeneration,
		Status: value.Status,
		Version: value.Version,
		CreatedAt: value.CreatedAt,
		UpdatedAt: value.UpdatedAt,
	}
}
