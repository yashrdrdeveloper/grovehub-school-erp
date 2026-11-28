// Package check rejects metadata-only directives and invalid annotation usage.
package check

import (
	"io/fs"
	"path/filepath"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/check/helpers"
)

func CheckServices(backend string) error {
	services := filepath.Join(backend, "services")
	providers := make(helpers.ProviderIndex)

	return filepath.WalkDir(services, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() {
			return helpers.CheckDirectory(entry)
		}

		if !helpers.IsMaintainedSource(path) {
			return nil
		}

		occurrences, err := helpers.CheckFile(path)
		if err != nil {
			return err
		}

		return providers.Add(occurrences)
	})
}

func CheckSource(path string) error {
	_, err := helpers.CheckFile(path)

	return err
}
