package helpers

import (
	"fmt"
	"strings"
)

func TrimMessageIDBrackets(value string) string {
	return strings.Trim(strings.TrimSpace(value), "<>")
}

func FormatMessageIDHeader(value string) string {
	return fmt.Sprintf("<%s>", TrimMessageIDBrackets(value))
}

func FormatReferenceMessageIDs(values []string) string {
	formatted := make([]string, 0, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			continue
		}
		formatted = append(formatted, FormatMessageIDHeader(value))
	}
	return strings.Join(formatted, " ")
}
