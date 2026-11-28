package project

import "grovehub/services/registry/internal/app/constants/status"

type CreateInput struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type StatusInput struct {
	ExpectedVersion int64         `json:"expectedVersion"`
	Status          status.Status `json:"status"`
}
