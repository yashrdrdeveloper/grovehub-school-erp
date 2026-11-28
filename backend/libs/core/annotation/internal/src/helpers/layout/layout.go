// Package layout resolves generator imports against the owning Go module.
package layout

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

// CoreImportPath names the current Core capability owner. Change this single
// generator contract together with the eventual standalone Core module move.
const CoreImportPath = "grovehub/libs/core"

func ReadModulePath(root string) (string, error) {
	filename := filepath.Join(root, "go.mod")
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("read generator target module: %w", err)
	}
	value := modfile.ModulePath(data)
	if value == "" {
		return "", fmt.Errorf("%s: missing module directive", filename)
	}
	if err := module.CheckImportPath(value); err != nil {
		return "", fmt.Errorf("%s: invalid module import path: %w", filename, err)
	}
	return value, nil
}

// ResolveImport only resolves packages inside the target module. Reject
// malformed or foreign paths instead of following a traversal outside root.
func ResolveImport(root string, modulePath string, importPath string) (string, error) {
	if err := module.CheckImportPath(importPath); err != nil {
		return "", err
	}
	prefix := modulePath + "/"
	if !strings.HasPrefix(importPath, prefix) {
		return "", fmt.Errorf("import %q is outside module %q", importPath, modulePath)
	}
	return ResolveRelative(root, strings.TrimPrefix(importPath, prefix))
}

// ResolveRelative joins a module-relative path without allowing an absolute
// path or a traversal outside the owning module.
func ResolveRelative(root, relative string) (string,error) {
	if !filepath.IsLocal(filepath.FromSlash(relative)) {
		return "", fmt.Errorf("invalid module-relative import %q", relative)
	}
	return filepath.Join(root, filepath.FromSlash(relative)), nil
}
