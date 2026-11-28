package helpers

import (
	"io/fs"
	"path/filepath"
	"strings"

	"grovehub/libs/core/annotation/internal/src/contracts"
	"grovehub/libs/core/annotation/internal/src/parser"
)

func CheckDirectory(entry fs.DirEntry) error {
	switch entry.Name() {
	case ".generated", "vendor":
		return filepath.SkipDir
	default:
		return nil
	}
}

func IsMaintainedSource(path string) bool {
	return strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go")
}

func CheckFile(path string) ([]contract.Occurrence, error) {
	occurrences, err := parser.ScanFile(path)
	if err != nil {
		return nil, err
	}

	for _, occurrence := range occurrences {
		if err := ValidateOccurrence(occurrence); err != nil {
			return nil, err
		}
	}

	if err := CheckValidationTags(path); err != nil {
		return nil, err
	}

	return occurrences, nil
}
