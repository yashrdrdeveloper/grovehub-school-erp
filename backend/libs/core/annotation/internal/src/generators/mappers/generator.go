// Package mappers generates strict typed DTO projection methods from @Mapper
// interfaces. Every destination field must have an explicit @Mapping rule.
// Business-specific transformations remain in maintained mapper methods.
package mappers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"grovehub/libs/core/annotation/internal/src/generators/mappers/helpers"
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
	for _, entry := range services {
		if !entry.IsDir() {
			continue
		}
		owner := filepath.Join(backend, "services", entry.Name())
		rootRel := "mappers"
		dir := filepath.Join(owner, rootRel)
		if stat, err := os.Stat(filepath.Join(owner, "internal", "app", "mappers")); err == nil && stat.IsDir() {
			rootRel = filepath.Join("internal", "app", "mappers")
			dir = filepath.Join(owner, rootRel)
		}
		if stat, err := os.Stat(dir); err != nil || !stat.IsDir() {
			continue
		}
		// Generate once per resource package, retaining legacy flat mappers.
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		directories := []string{"."}
		for _, entry := range entries {
			if entry.IsDir() && entry.Name() != "impl" && !strings.HasPrefix(entry.Name(), ".") {
				directories = append(directories, entry.Name())
			}
		}
		for _, rel := range directories {
			resourceDir := filepath.Join(dir, rel)
			group, found, err := helpers.Discover(backend, modulePath, resourceDir)
			if err != nil {
				return fmt.Errorf("%s mapper %s: %w", entry.Name(), resourceDir, err)
			}
			outDir := filepath.Join(owner, ".generated", rootRel)
			if rel != "." {
				outDir = filepath.Join(outDir, rel)
			}
			out := filepath.Join(outDir, "impl", "zz_mappers.gen.go")
			if !found || !helpers.HasMethods(group) {
				if err := os.Remove(out); err != nil && !os.IsNotExist(err) {
					return err
				}
				continue
			}
			if err := helpers.Write(group, out); err != nil {
				return fmt.Errorf("%s mapper %s: %w", entry.Name(), resourceDir, err)
			}
		}
	}
	return nil
}
