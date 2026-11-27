package helpers

import (
	"time"

	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/services/approval/operations/state"
)

func Snapshot(request entity.ApprovalRequest) state.Snapshot {
	return state.Snapshot{TenantID: request.TenantID, RequesterID: request.RequesterID, Status: request.Status, CurrentStep: request.CurrentStep, StepCount: request.StepCount, Version: request.Version, ExpiresAt: request.ExpiresAt}
}

func Advance(request entity.ApprovalRequest, next state.Snapshot, now time.Time) entity.ApprovalRequest {
	request.Status = next.Status
	request.CurrentStep = next.CurrentStep
	request.Version = next.Version
	request.UpdatedAt = now
	return request
}
