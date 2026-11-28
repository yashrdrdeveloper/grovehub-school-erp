package helpers

import (
	"fmt"
	"strings"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/contracts"
)

// ProviderIndex checks uniqueness across maintained source files.
type ProviderIndex map[string]string

func (index ProviderIndex) Add(occurrences []contract.Occurrence) error {
	for _, occurrence := range occurrences {
		if occurrence.Directive.Name != "Provider" {
			continue
		}

		id, _, err := ParseProviderDirective(occurrence.Directive.Arguments)
		if err != nil {
			return fmt.Errorf("%s:%d: %w", occurrence.File, occurrence.Line, err)
		}

		if earlier, exists := index[id]; exists {
			return fmt.Errorf("%s:%d: duplicate @Provider id=%s (first declared at %s)", occurrence.File, occurrence.Line, id, earlier)
		}

		index[id] = fmt.Sprintf("%s:%d", occurrence.File, occurrence.Line)
	}

	return nil
}

func ParseProviderDirective(args string) (string, string, error) {
	fields := strings.Fields(args)
	if len(fields) != 2 {
		return "", "", fmt.Errorf("@Provider requires id=NAME category=name")
	}

	values := make(map[string]string, 2)

	for _, field := range fields {
		key, value, ok := strings.Cut(field, "=")
		if !ok || value == "" || !isProviderKey(key) || values[key] != "" {
			return "", "", fmt.Errorf("invalid @Provider argument %q", field)
		}

		values[key] = value
	}

	id, category := values["id"], values["category"]
	if id == "" || category == "" || id != strings.ToUpper(id) {
		return "", "", fmt.Errorf("@Provider requires uppercase id and lowercase category")
	}

	if !validProviderID(id) {
		return "", "", fmt.Errorf("invalid @Provider id %q", id)
	}

	if !validProviderCategory(category) {
		return "", "", fmt.Errorf("invalid @Provider category %q", category)
	}

	return id, category, nil
}

func isProviderKey(key string) bool {
	return key == "id" || key == "category"
}

func validProviderID(id string) bool {
	for _, char := range id {
		if !(char >= 'A' && char <= 'Z') && !isDigit(char) && char != '_' {
			return false
		}
	}

	return true
}

func validProviderCategory(category string) bool {
	for _, char := range category {
		if !(char >= 'a' && char <= 'z') && !isDigit(char) && char != '_' && char != '-' {
			return false
		}
	}

	return true
}

func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}
