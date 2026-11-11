package errorcode

import (
	"net/http"

	coreerrors "grovehub/libs/core/errors"
)

var (
	InvalidPersonGID = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_INVALID_GID",
		"Person GID is invalid",
		http.StatusBadRequest,
	)
	InvalidPersonRequest = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_INVALID_REQUEST",
		"Person request is invalid",
		http.StatusBadRequest,
	)
	InvalidPersonContact = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_INVALID_CONTACT",
		"Person contact is invalid",
		http.StatusBadRequest,
	)
	InvalidPersonIdempotencyKey = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_IDEMPOTENCY_KEY_INVALID",
		"A valid idempotency key is required for this person create request",
		http.StatusBadRequest,
	)
	PersonIdempotencyConflict = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_IDEMPOTENCY_CONFLICT",
		"Idempotency key was already used for a different person create request",
		http.StatusConflict,
	)
	PersonContactConflict = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_CONTACT_CONFLICT",
		"Person contacts resolve to conflicting people",
		http.StatusConflict,
	)
	PersonConcurrentModification = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_CONCURRENT_MODIFICATION",
		"Person was modified by another request",
		http.StatusConflict,
	)
	InvalidPersonSearch = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_INVALID_SEARCH",
		"Person search request is invalid",
		http.StatusBadRequest,
	)
	InvalidPersonSearchCursor = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_INVALID_SEARCH_CURSOR",
		"Person search cursor is invalid",
		http.StatusBadRequest,
	)
	PersonNotFound = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_NOT_FOUND",
		"Person was not found",
		http.StatusNotFound,
	)
	PersonCreateFailed = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_CREATE_FAILED",
		"Unable to create person",
		http.StatusInternalServerError,
	)
	PersonReadFailed = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_READ_FAILED",
		"Unable to read person",
		http.StatusInternalServerError,
	)
	PersonSearchFailed = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_SEARCH_FAILED",
		"Unable to search people",
		http.StatusInternalServerError,
	)
	PersonUpdateFailed = coreerrors.NewDefinition(
		"FOUNDATION_PERSON_UPDATE_FAILED",
		"Unable to update person",
		http.StatusInternalServerError,
	)
)
