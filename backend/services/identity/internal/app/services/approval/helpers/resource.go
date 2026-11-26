package helpers

import (
	"encoding/json"

	kind "grovehub/services/identity/internal/app/constants/approval"
	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
	repo "grovehub/services/identity/internal/app/repository/approval"
)

func ResourceID(request entity.ApprovalRequest) (string, bool, error) {
	switch request.RequestType {
	case kind.RoleLifecycleChange, kind.AccountLifecycleChange:
		var intent dto.LifecycleIntent
		if err := json.Unmarshal(request.Payload, &intent); err != nil {
			return "", false, err
		}
		id, err := NormalizeID(intent.TargetID)
		return id, request.RequestType == kind.RoleLifecycleChange, err
	case kind.RolePermissionChange:
		var intent dto.PermissionIntent
		if err := json.Unmarshal(request.Payload, &intent); err != nil {
			return "", false, err
		}
		id, err := NormalizeID(intent.RoleID)
		return id, true, err
	case kind.RoleAssignmentChange:
		var intent dto.AssignmentIntent
		if err := json.Unmarshal(request.Payload, &intent); err != nil {
			return "", false, err
		}
		id, err := NormalizeID(intent.AccountID)
		return id, false, err
	default:
		return "", false, repo.ErrInvalid
	}
}
