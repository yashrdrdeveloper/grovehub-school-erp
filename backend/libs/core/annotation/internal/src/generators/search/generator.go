// Package search generates repository search plumbing from an explicit
// @Searchable declaration. SQL field expressions stay developer-owned in
// coresearch.Config; no query, join, or field is inferred from user input.
package search

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"grovehub/libs/core/annotation/internal/src/generators/search/helpers"
	"grovehub/libs/core/annotation/internal/src/helpers/layout"
)

func GenerateAll(backend string) error {
	modulePath, err := layout.ReadModulePath(backend)
	if err != nil {
		return err
	}
	services, err := os.ReadDir(filepath.Join(backend, "services"))
	if err != nil {
		return err
	}
	for _, service := range services {
		if !service.IsDir() {
			continue
		}
		owner := filepath.Join(backend, "services", service.Name())
		rootRel := "repository"
		dir := filepath.Join(owner, rootRel)
		if info, err := os.Stat(filepath.Join(owner, "internal", "app", "repository")); err == nil && info.IsDir() {
			rootRel = filepath.Join("internal", "app", "repository")
			dir = filepath.Join(owner, rootRel)
		}
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return err
		}
		// Search contracts are scoped to repository/<resource>, preserving
		// repository/impl as a legacy flat layout for other services.
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		directories := []string{"."}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			if entry.Name() == "impl" || entry.Name() == "constants" || entry.Name() == "helpers" || strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			directories = append(directories, entry.Name())
		}
		for _, rel := range directories {
			resourceDir := filepath.Join(dir, rel)
			specs, err := helpers.Discover(resourceDir, modulePath, service.Name())
			if err != nil {
				return fmt.Errorf("%s search %s: %w", service.Name(), resourceDir, err)
			}
			outDir := filepath.Join(owner, ".generated", rootRel)
			if rel != "." {
				outDir = filepath.Join(outDir, rel)
			}
			out := filepath.Join(outDir, "impl", "zz_search.gen.go")
			if len(specs) == 0 {
				if err := os.Remove(out); err != nil && !os.IsNotExist(err) {
					return err
				}
				continue
			}
			if err := helpers.ValidateImplementation(filepath.Join(resourceDir, "impl"), specs); err != nil {
				return err
			}
			if err := helpers.Write(out, modulePath, specs); err != nil {
				return err
			}
		}
	}
	return nil
}
