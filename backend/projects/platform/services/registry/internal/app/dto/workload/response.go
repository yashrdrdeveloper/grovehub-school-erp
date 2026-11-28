package workload

import (
	"time"

	"grovehub/services/registry/internal/app/constants/status"
	"grovehub/services/registry/internal/app/entity"
)

type Response struct {
	WorkloadID string        `json:"workloadId"`
	ProjectID  *string       `json:"projectId,omitempty"`
	Subject    string        `json:"subject"`
	Issuer     string        `json:"issuer"`
	Status     status.Status `json:"status"`
	Version    int64         `json:"version"`
	CreatedAt  time.Time     `json:"createdAt"`
	UpdatedAt  time.Time     `json:"updatedAt"`
}

func FromEntity(value entity.Workload) Response {
	return Response{
		WorkloadID: value.WorkloadID,
		ProjectID: value.ProjectID,
		Subject: value.Subject,
		Issuer: value.Issuer,
		Status: value.Status,
		Version: value.Version,
		CreatedAt: value.CreatedAt,
		UpdatedAt: value.UpdatedAt,
	}
}
