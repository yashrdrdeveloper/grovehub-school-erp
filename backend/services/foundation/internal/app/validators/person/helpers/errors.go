package helpers

import (
    "errors"
    "strings"

    annotationvalidation "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/validation"
    coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
    personerror "grovehub/services/foundation/internal/app/constants/error"
)

func invalidPersonSearch() error {
	return coreerrors.NewError(personerror.InvalidPersonSearch)
}

func invalidPersonSearchCursor() error {
	return coreerrors.NewError(personerror.InvalidPersonSearchCursor)
}

func MapPersonDTOError(err error) error {
	var fieldError annotationvalidation.Error
	if !errors.As(err, &fieldError) { return err }
	field := strings.TrimPrefix(fieldError.Field, "RequestBody.")
	if strings.HasPrefix(field, "Contacts.") || field == "Contacts" {
		return coreerrors.NewError(personerror.InvalidPersonContact)
	}
	return coreerrors.NewError(personerror.InvalidPersonRequest)
}
