// Package clients generates typed client methods from @Client contracts.
//
// @RequestLine calls the verified generated controller route in-process,
// preserving the HTTP binding, validation and response pipeline.
// Legacy @Proxy methods remain supported for other service clients.
package clients

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"grovehub/libs/core/annotation/internal/src/generators/clients/helpers"
)

func GenerateAll(backend string) error {
	services, err := os.ReadDir(filepath.Join(backend, "services"))
	if err != nil { return err }
	for _,service := range services {
		if !service.IsDir() { continue }
		owner := filepath.Join(backend,"services",service.Name())
		clientDir := filepath.Join(owner,"client")
		if info,err := os.Stat(clientDir); err != nil || !info.IsDir() { continue }
		// Discover both legacy flat clients and client/<resource> contracts.
		entries,err:=os.ReadDir(clientDir)
		if err!=nil {return err}
		directories:=[]string{"."}
		for _,entry:=range entries {
			if entry.IsDir() && entry.Name()!="impl" && !strings.HasPrefix(entry.Name(),".") {
				directories=append(directories,entry.Name())
			}
		}
		for _,rel:=range directories {
			dir:=filepath.Join(clientDir,rel)
			generatedRoutes := filepath.Join(owner,".generated","routes","zz_routes.gen.go")
            if _, err := os.Stat(filepath.Join(owner, "internal", "app", "routes")); err == nil {
                generatedRoutes = filepath.Join(owner, ".generated", "internal", "app", "routes", "zz_routes.gen.go")
            }
            proxies,err:=helpers.Discover(dir,generatedRoutes)
			if err!=nil{return fmt.Errorf("%s %s: %w",service.Name(),dir,err)}
			outDir:=filepath.Join(owner,".generated","client")
            if rel!="." {outDir=filepath.Join(outDir,rel)}
            if _, err := os.Stat(filepath.Join(owner,"internal","app","client",rel,"impl")); err == nil {
                outDir = filepath.Join(owner,".generated","internal","app","client",rel)
            }
			output:=filepath.Join(outDir,"impl","zz_clients.gen.go")
			if len(proxies)==0 {
				if err:=os.Remove(output);err!=nil&&!os.IsNotExist(err){return err}
				continue
			}
			if err:=helpers.WriteProxies(proxies,output);err!=nil{return fmt.Errorf("%s %s: %w",service.Name(),dir,err)}
		}
	}
	return nil
}
