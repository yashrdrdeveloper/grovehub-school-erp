// Package entities generates optional builders/getters/setters only when
// explicitly requested on NON-entity types. @Entity itself only maps a GORM
// table and must not imply any constructor or accessor generation.
// Code is written only to the owning service's ignored .generated/entity tree.
package entities

import (
	"fmt"
	"os"
	"path/filepath"

	"grovehub/libs/core/annotation/internal/src/generators/entities/helpers"
)

func GenerateAll(backend string) error {
	services, err := os.ReadDir(filepath.Join(backend, "services"))
	if err != nil {
		return err
	}
	for _, service := range services {
		if !service.IsDir() {
			continue
		}
		owner := filepath.Join(backend, "services", service.Name())
		rel := "entity"
		dir := filepath.Join(owner, rel)
		if info, err := os.Stat(filepath.Join(owner, "internal", "app", "entity")); err == nil && info.IsDir() {
			rel = filepath.Join("internal", "app", "entity")
			dir = filepath.Join(owner, rel)
		}
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			continue
		}
		if err := helpers.Generate(dir, filepath.Join(owner, ".generated", rel)); err != nil {
			return fmt.Errorf("%s entities: %w", service.Name(), err)
		}
	}
	return nil
}
