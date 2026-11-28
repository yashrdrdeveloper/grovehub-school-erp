package pack

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"golang.org/x/mod/zip"
)

// Run prepares a versioned module in the standard Go proxy layout.
func Run(args []string) error {
	if len(args) != 4 {
		return fmt.Errorf("usage: modulepack SOURCE PROXY VERSION RFC3339_TIME")
	}
	source, proxy, version := args[0], args[1], args[2]
	mod, err := os.ReadFile(filepath.Join(source, "go.mod"))
	if err != nil {
		return err
	}
	parsed, err := modfile.Parse("go.mod", mod, nil)
	if err != nil {
		return err
	}
	if parsed.Module == nil || len(parsed.Replace) != 0 {
		return fmt.Errorf("module must declare its path and have no replacements")
	}
	m := module.Version{Path: parsed.Module.Mod.Path, Version: version}
	if err := module.Check(m.Path, m.Version); err != nil {
		return err
	}
	escapedPath, err := module.EscapePath(m.Path)
	if err != nil {
		return err
	}
	escapedVersion, err := module.EscapeVersion(m.Version)
	if err != nil {
		return err
	}
	stamp, err := time.Parse(time.RFC3339, args[3])
	if err != nil {
		return err
	}
	dir := filepath.Join(proxy, escapedPath, "@v")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	base := filepath.Join(dir, escapedVersion)
	archive, err := os.OpenFile(base+".zip", os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	err = zip.CreateFromDir(archive, m, source)
	closeErr := archive.Close()
	if err != nil || closeErr != nil {
		os.Remove(base + ".zip")
		if err != nil {
			return err
		}
		return closeErr
	}
	info, err := json.Marshal(struct {
		Version string
		Time time.Time
	}{version, stamp.UTC()})
	if err != nil {
		return err
	}
	if err := os.WriteFile(base+".mod", mod, 0644); err != nil {
		return err
	}
	if err := os.WriteFile(base+".info", info, 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "list"), []byte(version+"\n"), 0644)
}

