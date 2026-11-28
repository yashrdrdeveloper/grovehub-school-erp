package project

import (
	"time"

	"grovehub/services/registry/internal/app/constants/status"
	"grovehub/services/registry/internal/app/entity"
)

type Response struct {
	ProjectID string `json:"projectId"`
	Code string `json:"code"`
	Name string `json:"name"`
	Status status.Status `json:"status"`
	Version int64 `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func FromEntity(value entity.Project) Response {
	return Response{
		ProjectID: value.ProjectID,
		Code: value.Code,
		Name: value.Name,
		Status: value.Status,
		Version: value.Version,
		CreatedAt: value.CreatedAt,
		UpdatedAt: value.UpdatedAt,
	}
}
