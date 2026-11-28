// Package patchdto generates JSON marshal methods for @PatchRequestDto
// structs into the existing service-local, ignored Go source overlay.
package patchdto

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/patchdto/helpers"
)

// DTO packages can live in legacy dto/, client/, or the owning service's
// internal/app/dto/. Generated code is overlaid into the original package.
func GenerateAll(backend string) error {
	services, err := os.ReadDir(filepath.Join(backend, "services"))
	if err != nil { return err }
	for _, entry := range services {
		if !entry.IsDir() { continue }
		owner := filepath.Join(backend, "services", entry.Name())
		roots := []string{filepath.Join(owner, "dto"), filepath.Join(owner, "client"), filepath.Join(owner, "internal", "app", "dto")}
        for _,root := range roots {
            if _, err := os.Stat(root); os.IsNotExist(err) { continue } else if err != nil { return err }
            err = filepath.WalkDir(root, func(dir string, item os.DirEntry, walkErr error) error {
			if walkErr != nil { return walkErr }
			if !item.IsDir() { return nil }
			if item.Name() == ".generated" || item.Name() == "vendor" { return filepath.SkipDir }
			specs, pkg, err := helpers.Discover(dir)
			if err != nil { return fmt.Errorf("%s: %w", dir, err) }
			rel, err := filepath.Rel(owner, dir)
			if err != nil { return err }
			output := filepath.Join(owner, ".generated", rel, "zz_patchdto.gen.go")
			if len(specs) == 0 {
				if err := os.Remove(output); err != nil && !os.IsNotExist(err) { return err }
				return nil
			}
			return helpers.Write(output, pkg, specs)
		})
		if err != nil { return err }
        }
	}
	return nil
}
