package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
	InvalidPersonGID = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_INVALID_GID",
		"Person GID is invalid",
	)
	InvalidPersonRequest = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_INVALID_REQUEST",
		"Person request is invalid",
	)
	InvalidPersonContact = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_INVALID_CONTACT",
		"Person contact is invalid",
	)
	InvalidPersonIdempotencyKey = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_IDEMPOTENCY_KEY_INVALID",
		"A valid idempotency key is required for this person create request",
	)
	PersonIdempotencyConflict = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_IDEMPOTENCY_CONFLICT",
		"Idempotency key was already used for a different person create request",
	)
	PersonContactConflict = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_CONTACT_CONFLICT",
		"Person contacts resolve to conflicting people",
	)
	PersonConcurrentModification = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_CONCURRENT_MODIFICATION",
		"Person was modified by another request",
	)
	InvalidPersonSearch = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_INVALID_SEARCH",
		"Person search request is invalid",
	)
	InvalidPersonSearchCursor = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_INVALID_SEARCH_CURSOR",
		"Person search cursor is invalid",
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
	PersonSearchFailed = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_SEARCH_FAILED",
		"Unable to search people",
	)
	PersonUpdateFailed = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_UPDATE_FAILED",
		"Unable to update person",
	)
)
