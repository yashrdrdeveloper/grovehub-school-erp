package helpers

import (
	"fmt"
	"go/ast"
	"strings"
)

func prefixFor(doc *ast.CommentGroup) (string, bool, error) {
	if doc == nil {
		return "", false, nil
	}
	found := false
	var prefix string
	for _, comment := range doc.List {
		line := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		if line != "@Routes" && !strings.HasPrefix(line, "@Routes ") {
			continue
		}
		if found {
			return "", true, fmt.Errorf("duplicate @Routes declaration")
		}
		found = true
		prefix = strings.TrimSpace(strings.TrimPrefix(line, "@Routes"))
		if err := validatePath(prefix); err != nil {
			return "", true, err
		}
	}
	return prefix, found, nil
}

func validatePath(value string) error {
	if !strings.HasPrefix(value, "/") || strings.Contains(value, "//") ||
		strings.ContainsAny(value, " \t\r\n?#") {
		return fmt.Errorf("invalid route path %q", value)
	}
	for _, part := range strings.Split(value, "/") {
		if part == "." || part == ".." {
			return fmt.Errorf("invalid route path %q", value)
		}
	}
	return nil
}

func joinPaths(parts ...string) string {
	segments := make([]string, 0, len(parts))
	for _, part := range parts {
		if segment := strings.Trim(part, "/"); segment != "" {
			segments = append(segments, segment)
		}
	}
	return "/" + strings.Join(segments, "/")
}
