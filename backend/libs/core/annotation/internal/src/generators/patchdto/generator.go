// Package patchdto generates JSON marshal methods for @PatchRequestDto
// structs into the existing service-local, ignored Go source overlay.
package patchdto

import (
	"grovehub/libs/core/annotation/internal/src/helpers/layout"
	"fmt"
	"os"
	"path/filepath"

	"grovehub/libs/core/annotation/internal/src/generators/patchdto/helpers"
)

// DTO packages can live in legacy dto/, client/, or the owning service's
// internal/app/dto/. Generated code is overlaid into the original package.
func GenerateAll(backend string) error {
	services, err := os.ReadDir(filepath.Join(backend, "services"))
	if err != nil { return err }
	for _, entry := range services {
		if !entry.IsDir() { continue }
		owner := filepath.Join(backend, "services", entry.Name())
		roots := []string{"dto","client",filepath.Join("internal","app","dto")}
        for _,root := range roots {
            err = layout.WalkOwnedDirectories(owner,root,map[string]bool{".generated":true,"vendor":true},func(dir,rel string)error{
			specs, pkg, err := helpers.Discover(dir)
			if err != nil { return fmt.Errorf("%s: %w", dir, err) }
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
