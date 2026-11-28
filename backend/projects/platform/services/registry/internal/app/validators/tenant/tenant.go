package tenant

import (
	"math"
	
	coreerrors "grovehub/libs/core/errors"
	"grovehub/libs/core/utils"

	definitions "grovehub/services/registry/internal/app/constants/error"
	"grovehub/services/registry/internal/app/constants/status"
	"grovehub/services/registry/internal/app/constants/validation"
	dto "grovehub/services/registry/internal/app/dto/tenant"
)

func Validate(projectID, tenantID string) error {
	if utils.ValidateUUID(projectID) != nil || utils.ValidateUUID(tenantID) != nil {
		return coreerrors.NewError(definitions.Invalid)
	}
	return nil
}

func ValidateCreate(projectID string, input dto.CreateInput) error {
	reference, err := utils.NormalizeText(input.ExternalReference, validation.MaxReferenceBytes)
	if utils.ValidateUUID(projectID) != nil || err != nil || reference == "" || reference != input.ExternalReference {
		return coreerrors.NewError(definitions.Invalid)
	}
	return nil
}

func ValidateStatus(projectID, tenantID string, input dto.StatusInput) error {
	if err := Validate(projectID, tenantID); err != nil {
		return err
	}
	if input.ExpectedVersion < 1 || input.ExpectedVersion == math.MaxInt64 || (input.Status != status.Active && input.Status != status.Disabled) {
		return coreerrors.NewError(definitions.Invalid)
	}
	return nil
}
