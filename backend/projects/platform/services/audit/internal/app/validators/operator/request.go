package operator

import (
	coreerrors "grovehub/libs/core/errors"
	definitions "grovehub/services/audit/internal/app/constants/error"
	dto "grovehub/services/audit/internal/app/dto/operator"
)

func ValidateWrite(input dto.WriteRequest) error {
	return validate(ValidActivity(input.RequestBody), definitions.OperatorInvalid)
}

func ValidateRename(input dto.RenameReceiptRequest) error {
	return validate(ValidRename(input.RequestBody), definitions.OperatorInvalid)
}

func ValidateStatus(input dto.StatusReceiptRequest) error {
	return validate(ValidStatus(input.RequestBody), definitions.OperatorInvalid)
}

func ValidateSearch(input dto.SearchRequest) error {
	return validate(ValidSearch(input.RequestBody), definitions.OperatorInvalid)
}

func ValidateReceiptSearch(input dto.ReceiptSearchRequest) error {
	return validate(ValidReceiptSearch(input.RequestBody), definitions.OperatorInvalid)
}

func ValidateSchoolCheck(input dto.SchoolCheckRequest) error {
	return validate(ValidID(input.Params.SchoolID), definitions.OperatorForbidden)
}

func validate(valid bool, definition coreerrors.Definition) error {
	if !valid {
		return coreerrors.NewError(definition)
	}
	return nil
}
