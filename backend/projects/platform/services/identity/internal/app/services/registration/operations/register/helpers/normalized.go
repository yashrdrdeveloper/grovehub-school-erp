package helpers

import (
	"strings"
	"unicode/utf8"

	idempotency "grovehub/libs/core/idempotency"
	dto "grovehub/services/identity/internal/app/dto/account"
	repo "grovehub/services/identity/internal/app/repository/account"
	validator "grovehub/services/identity/internal/app/validators/account"
)

func NormalizePending(request dto.Create, key string, normalizer validator.Validator) (dto.Create, string, string, error) {
	// Initial assignment is an independent approved operation. Never silently
	// discard a caller's requested roles or write legacy role columns.
	if request.RoleAssignment != nil || normalizer == nil {
		return dto.Create{}, "", "", repo.ErrConflict
	}
	normalized, err := normalizer.NormalizeCreate(request)
	if err != nil {
		return dto.Create{}, "", "", err
	}
	values := []string{normalized.Person.FirstName, normalized.Person.LastName}
	if normalized.Person.MiddleName != nil {
		values = append(values, *normalized.Person.MiddleName)
	}
	for _, value := range values {
		if !utf8.ValidString(value) || strings.ContainsRune(value, 0) || len(value) > 100 {
			return dto.Create{}, "", "", repo.ErrConflict
		}
	}
	if strings.TrimSpace(normalized.PhoneNumber) == "" {
		normalized.PhoneNumber = ""
	}
	key, err = idempotency.ValidateKey(key)
	if err != nil {
		return dto.Create{}, "", "", err
	}
	digest, err := idempotency.Fingerprint(normalized)
	return normalized, key, digest, err
}
