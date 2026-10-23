package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
	InvalidPersonGID = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_INVALID_GID",
		"Person GID is invalid",
	)
	PersonNotFound = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_NOT_FOUND",
		"Person was not found",
	)
	PersonCreateFailed = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_CREATE_FAILED",
		"Unable to create person",
	)
	PersonReadFailed = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_READ_FAILED",
		"Unable to read person",
	)
	PersonUpdateFailed = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_UPDATE_FAILED",
		"Unable to update person",
	)
)
