package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/permission"
	dto "grovehub/services/identity/internal/app/dto/approval"
	repo "grovehub/services/identity/internal/app/repository/approval"
)

func NormalizeID(value string) (string, error) {
	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil {
		return "", repo.ErrInvalid
	}
	return id.String(), nil
}

func NormalizeLifecycle(request dto.CreateLifecycle) (dto.CreateLifecycle, permission.Code, []byte, string, error) {
	target, err := NormalizeID(request.Intent.TargetID)
	if err != nil || request.Intent.Version < 1 {
		return dto.CreateLifecycle{}, "", nil, "", repo.ErrInvalid
	}
	request.Intent.TargetID = target
	request.Reason = strings.TrimSpace(request.Reason)
	if !validText(request.Reason, 2000) || !validText(request.IdempotencyKey, 128) || strings.TrimSpace(request.IdempotencyKey) != request.IdempotencyKey {
		return dto.CreateLifecycle{}, "", nil, "", repo.ErrInvalid
	}
	code, known := LifecyclePermission(request.Type, request.Intent.Action)
	if !known {
		return dto.CreateLifecycle{}, "", nil, "", repo.ErrInvalid
	}
	payload, err := json.Marshal(request.Intent)
	if err != nil {
		return dto.CreateLifecycle{}, "", nil, "", err
	}
	digest := sha256.Sum256(payload)
	return request, code, payload, hex.EncodeToString(digest[:]), nil
}

func LifecyclePermission(kind approval.RequestType, action string) (permission.Code, bool) {
	switch kind {
	case approval.AccountLifecycleChange:
		switch action {
		case "suspend":
			return permission.AccountSuspend, true
		case "activate":
			return permission.AccountActivate, true
		case "deactivate":
			return permission.AccountDeactivate, true
		}
	case approval.RoleLifecycleChange:
		switch action {
		case "activate":
			return permission.RoleUpdate, true
		case "deactivate":
			return permission.RoleDeactivate, true
		}
	}
	return "", false
}

func validText(value string, limit int) bool {
	return strings.TrimSpace(value) != "" && utf8.ValidString(value) && !strings.ContainsRune(value, 0) && utf8.RuneCountInString(value) <= limit
}

func ValidatePolicy(policy repo.Policy) error {
	if _, err := NormalizeID(policy.PolicyID); err != nil || len(policy.Steps) < 1 || len(policy.Steps) > approval.MaxSteps {
		return repo.ErrPolicy
	}
	for index, step := range policy.Steps {
		if step.Number != index+1 {
			return repo.ErrPolicy
		}
		if _, known := permission.Find(step.Permission); !known || !permission.SchoolDelegable(step.Permission) {
			return repo.ErrPolicy
		}
		if step.RoleID != nil {
			if _, err := NormalizeID(*step.RoleID); err != nil {
				return repo.ErrPolicy
			}
		}
	}
	return nil
}
