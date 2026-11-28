package helpers

import (
	"bytes"
	"encoding/json"
	"io"

	kind "grovehub/services/identity/internal/app/constants/approval"
	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
	repo "grovehub/services/identity/internal/app/repository/approval"
)

// VerifyIntent decodes the typed persisted JSONB object and recomputes the
// canonical digest. PostgreSQL JSONB may change key order or whitespace.
func VerifyIntent(request entity.ApprovalRequest) error {
	var digest string
	switch request.RequestType {
	case kind.AccountLifecycleChange, kind.RoleLifecycleChange:
		var intent dto.LifecycleIntent
		if err := decode(request.Payload, &intent); err != nil {
			return err
		}
		_, _, _, actual, err := NormalizeLifecycle(dto.CreateLifecycle{Type: request.RequestType, Intent: intent, Reason: request.Reason, IdempotencyKey: request.IdempotencyKey})
		if err != nil {
			return err
		}
		digest = actual
	case kind.RolePermissionChange:
		var intent dto.PermissionIntent
		if err := decode(request.Payload, &intent); err != nil {
			return err
		}
		normalized, err := NormalizePermissions(intent)
		if err != nil {
			return err
		}
		_, digest, err = EncodeIntent(normalized)
		if err != nil {
			return err
		}
	case kind.RoleAssignmentChange:
		var intent dto.AssignmentIntent
		if err := decode(request.Payload, &intent); err != nil {
			return err
		}
		normalized, err := NormalizeAssignment(intent)
		if err != nil {
			return err
		}
		_, digest, err = EncodeIntent(normalized)
		if err != nil {
			return err
		}
	default:
		return repo.ErrInvalid
	}
	if digest != request.PayloadSHA256 {
		return repo.ErrConflict
	}
	return nil
}

func decode(payload []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return repo.ErrInvalid
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return repo.ErrInvalid
	}
	return nil
}
