package workload

import (
	"crypto/ed25519"
	"encoding/base64"
	"math"
	"regexp"
	"time"

	coreerrors "grovehub/libs/core/errors"
	"grovehub/libs/core/utils"
	"grovehub/libs/platform/scope"

	definitions "grovehub/services/registry/internal/app/constants/error"
	"grovehub/services/registry/internal/app/constants/status"
	dto "grovehub/services/registry/internal/app/dto/workload"
)

const MaxKeyLifetime = 366 * 24 * time.Hour

var identifier = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9._:-]{0,127}$")
var keyIdentifier = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_:-]{0,127}$")

func ValidateKeyID(value string) error {
	if !keyIdentifier.MatchString(value) {
		return coreerrors.NewError(definitions.Invalid)
	}
	return nil
}

func ValidateID(value string) error {
	if utils.ValidateUUID(value) != nil {
		return coreerrors.NewError(definitions.Invalid)
	}
	return nil
}

func ValidateIdentifier(value string) error {
	if !identifier.MatchString(value) {
		return coreerrors.NewError(definitions.Invalid)
	}
	return nil
}

func ValidateCreate(input dto.CreateInput) error {
	if (input.ProjectID != nil && ValidateID(*input.ProjectID) != nil) ||
		ValidateIdentifier(input.Issuer) != nil || ValidateIdentifier(input.Subject) != nil {
		return coreerrors.NewError(definitions.Invalid)
	}
	return nil
}

func ValidateStatus(input dto.StatusInput) error {
	if input.ExpectedVersion < 1 || input.ExpectedVersion == math.MaxInt64 ||
		(input.Status != status.Active && input.Status != status.Disabled) {
		return coreerrors.NewError(definitions.Invalid)
	}
	return nil
}

func ValidateKey(input dto.KeyInput, now time.Time) ([]byte, error) {
	key, err := base64.RawURLEncoding.Strict().DecodeString(input.PublicKey)
	if ValidateKeyID(input.KeyID) != nil || err != nil ||
		len(key) != ed25519.PublicKeySize || base64.RawURLEncoding.EncodeToString(key) != input.PublicKey ||
		input.ValidFrom.IsZero() || !input.ValidUntil.After(now) ||
		!input.ValidUntil.After(input.ValidFrom) || input.ValidUntil.Sub(input.ValidFrom) > MaxKeyLifetime {
		return nil, coreerrors.NewError(definitions.Invalid)
	}
	return key, nil
}

func ValidateGrant(input dto.GrantInput) error {
	target := scope.Target{Level: input.Scope, ProjectID: input.TargetProjectID, TenantID: input.TargetTenantID}
	if scope.ValidateTarget(target) != nil || ValidateIdentifier(input.Capability) != nil ||
		(input.TargetProjectID != "" && ValidateID(input.TargetProjectID) != nil) ||
		(input.TargetTenantID != "" && ValidateID(input.TargetTenantID) != nil) ||
		(input.Scope == scope.Platform && input.AllowDelegation) {
		return coreerrors.NewError(definitions.Invalid)
	}
	return nil
}
