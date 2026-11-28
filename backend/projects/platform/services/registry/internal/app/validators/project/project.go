package project

import (
	"math"
	"regexp"
	
	coreerrors "grovehub/libs/core/errors"
	"grovehub/libs/core/utils"

	definitions "grovehub/services/registry/internal/app/constants/error"
	"grovehub/services/registry/internal/app/constants/status"
	"grovehub/services/registry/internal/app/constants/validation"
	dto "grovehub/services/registry/internal/app/dto/project"
)

func Validate(projectID string) error {
	if utils.ValidateUUID(projectID) != nil {
		return coreerrors.NewError(definitions.Invalid)
	}
	return nil
}

var codePattern = regexp.MustCompile(validation.ProjectCodePattern)

func ValidateCreate(input dto.CreateInput) error {
	name, err := utils.NormalizeText(input.Name, validation.MaxNameBytes)
	if !codePattern.MatchString(input.Code) || err != nil || name == "" || name != input.Name {
		return coreerrors.NewError(definitions.Invalid)
	}
	return nil
}

func ValidateStatus(projectID string, input dto.StatusInput) error {
	if err := Validate(projectID); err != nil {
		return err
	}
	if input.ExpectedVersion < 1 || input.ExpectedVersion == math.MaxInt64 || (input.Status != status.Active && input.Status != status.Disabled) {
		return coreerrors.NewError(definitions.Invalid)
	}
	return nil
}
