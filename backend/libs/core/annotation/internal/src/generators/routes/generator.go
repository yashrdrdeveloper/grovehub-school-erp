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

	"grovehub/libs/core/annotation/internal/src/generators/routes/helpers"
	"grovehub/libs/core/annotation/internal/src/helpers/layout"
)

func GenerateAll(backend string) error {
	modulePath, err := layout.ReadModulePath(backend)
	if err != nil {
		return err
	}
	_, markerErr := os.Stat(filepath.Join(backend, ".restricted-routes"))
	if markerErr != nil && !os.IsNotExist(markerErr) {
		return markerErr
	}
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
		controllerRel := "controller"
		controllerDir := filepath.Join(owner, controllerRel)
		if info, err := os.Stat(filepath.Join(owner, "internal", "app", "controller")); err == nil && info.IsDir() {
			controllerRel = filepath.Join("internal", "app", "controller")
			controllerDir = filepath.Join(owner, controllerRel)
		}
		if info, err := os.Stat(controllerDir); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return err
		} else if !info.IsDir() {
			return fmt.Errorf("%s: controller must be a directory", name)
		}

		routesRel := filepath.Join(filepath.Dir(controllerRel), "routes")
		routesDir := filepath.Join(owner, routesRel)
		servicePath, configured, err := helpers.ServicePrefix(filepath.Join(routesDir, "http.go"))
		if err != nil {
			return fmt.Errorf("%s routes: %w", name, err)
		}
		groups, err := helpers.DiscoverWithTier(controllerDir, servicePath, configured, restrictedOnly)
		if err != nil {
			return fmt.Errorf("%s routes: %w", name, err)
		}
		output := filepath.Join(owner, ".generated", routesRel, "zz_routes.gen.go")
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
	return generateLauncher(backend, modulePath)
}

// Launcher owns only the two explicit infrastructure probes. It cannot be
// used to expose service controllers outside their restricted route tiers.
func generateLauncher(backend, modulePath string) error {
	owner := filepath.Join(backend, "launcher")
	controller := filepath.Join(owner, "internal", "src", "controller")
	if _, err := os.Stat(controller); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	groups, err := helpers.DiscoverWithTier(controller, "", false, false)
	if err != nil {
		return err
	}
	return helpers.WriteLauncherRoutes(modulePath, groups, filepath.Join(owner, ".generated", "internal", "src", "routes", "zz_routes.gen.go"))
}
