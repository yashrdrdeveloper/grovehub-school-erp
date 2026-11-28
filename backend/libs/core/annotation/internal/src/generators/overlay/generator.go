// Package overlay maps generated files in each module's .generated subtree
// into their maintained Go package locations without writing files there.
package overlay

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/layout"
	generatedoutput "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/output"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Configuration struct {
	Replace map[string]string `json:"Replace"`
}



func WriteOverlay(backend, core string) error {
	root, err := filepath.Abs(backend)
	if err != nil {
		return err
	}
	core, err = filepath.Abs(core)
	if err != nil {
		return err
	}
	replacement := map[string]string{}

	scan := func(owner string) error {
		return layout.WalkOwned(owner,".generated",func(path,relative string,item fs.DirEntry)error{
			if item.IsDir() || !strings.HasSuffix(item.Name(), ".go") {
				return nil
			}
			if !item.Type().IsRegular() {
				return fmt.Errorf("generated Go file %s must be regular", path)
			}
			rel:=strings.TrimPrefix(relative, ".generated"+string(os.PathSeparator))
			logical := filepath.Join(owner, rel)
			if _, err := os.Stat(filepath.Dir(logical)); err != nil {
				return fmt.Errorf("generated source has no original package directory %s: %w", logical, err)
			}
			if _, err := os.Stat(logical); err == nil {
				return fmt.Errorf("generated code conflicts with maintained source: %s", logical)
			} else if !os.IsNotExist(err) {
				return err
			}
			if prior, ok := replacement[logical]; ok {
				return fmt.Errorf("duplicate virtual file %s: %s and %s", logical, prior, path)
			}
			replacement[logical] = path
			return nil
		})
	}

	if root != core {
		services, err := os.ReadDir(filepath.Join(root, "services"))
		if err != nil {
			return err
		}
		for _, service := range services {
			if !service.IsDir() {
				continue
			}
			if err := scan(filepath.Join(root, "services", service.Name())); err != nil {
				return err
			}
		}
	}
	if err := scan(core); err != nil {
		return err
	}
	if len(replacement) == 0 && root != core {
		return fmt.Errorf("no generated Go files found; run generate-annotations.sh first")
	}

	return generatedoutput.WriteJSON(filepath.Join(root, ".generated", "overlay.json"), "overlay-*.json", Configuration{Replace:replacement})
}
