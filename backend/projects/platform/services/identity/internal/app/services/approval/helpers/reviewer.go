package helpers

import (
	kind "grovehub/services/identity/internal/app/constants/approval"
	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/services/approval/operations/state"
)

func ReviewerIndependent(request entity.ApprovalRequest, reviewerID string) error {
	if reviewerID == "" || reviewerID == request.RequesterID {
		return state.ErrForbidden
	}
	if request.RequestType == kind.RoleAssignmentChange {
		var intent dto.AssignmentIntent
		if err := decode(request.Payload, &intent); err != nil {
			return err
		}
		target, err := NormalizeID(intent.AccountID)
		if err != nil || target == reviewerID {
			return state.ErrForbidden
		}
	}
	return nil
}
