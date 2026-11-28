// Package routes generates typed net/http registrations from annotated
// controller interfaces. The service-level @Routes declaration lives on
// routes.RegisterHTTP; controller package paths declare access tiers.
//
// This generator only groups routes. No tier enforces authentication,
// authorization or internal-network restrictions.
package routes

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/routes/helpers"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/layout"
)

func GenerateAll(backend string) error {
	modulePath, err := layout.ReadModulePath(backend)
	if err != nil { return err }
	_, markerErr := os.Stat(filepath.Join(backend, ".restricted-routes"))
	if markerErr != nil && !os.IsNotExist(markerErr) { return markerErr }
	restrictedOnly := markerErr == nil
	services, err := os.ReadDir(filepath.Join(backend, "services"))
	if err != nil {
		return err
	}
	for _, service := range services {
		if !service.IsDir() {
			continue
		}
		name := service.Name()
		owner := filepath.Join(backend, "services", name)
		controllerDir := filepath.Join(owner, "controller")
        if info, err := os.Stat(filepath.Join(owner, "internal", "app", "controller")); err == nil && info.IsDir() {
            controllerDir = filepath.Join(owner, "internal", "app", "controller")
        }
		if info, err := os.Stat(controllerDir); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return err
		} else if !info.IsDir() {
			return fmt.Errorf("%s: controller must be a directory", name)
		}

		routesDir := filepath.Join(filepath.Dir(controllerDir), "routes")
        servicePath, configured, err := helpers.ServicePrefix(filepath.Join(routesDir, "http.go"))
		if err != nil {
			return fmt.Errorf("%s routes: %w", name, err)
		}
		groups, err := helpers.DiscoverWithTier(controllerDir, servicePath, configured, restrictedOnly)
		if err != nil {
			return fmt.Errorf("%s routes: %w", name, err)
		}
		rel, err := filepath.Rel(owner, routesDir)
        if err != nil { return err }
        output := filepath.Join(owner, ".generated", rel, "zz_routes.gen.go")
		if len(groups) == 0 {
			if err := os.Remove(output); err != nil && !os.IsNotExist(err) {
				return err
			}
			continue
		}
		if err := helpers.WriteRoutes(modulePath, name, groups, output); err != nil {
			return fmt.Errorf("%s routes: %w", name, err)
		}
	}
	return nil
}
