package helpers

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/layout"
	"fmt"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

// resolveConfig finds an explicitly imported repository search configuration
// in a subfolder by its declared Go package, not by an inferred SQL/table name.
func resolveConfig(repositoryDir,modulePath,service,reference string)(dir,importPath,alias,variable string,err error) {
	parts:=strings.Split(reference,".")
	if len(parts)==1 {
		if !token.IsIdentifier(parts[0]) {return "","","","",fmt.Errorf("invalid @Searchable config %q",reference)}
		return filepath.Join(repositoryDir,"impl"),"","",parts[0],nil
	}
	if len(parts)!=2||!token.IsIdentifier(parts[0])||!token.IsIdentifier(parts[1]) {
		return "","","","",fmt.Errorf("invalid @Searchable config %q",reference)
	}
	// Repository-owned search configs may live in repository/search (legacy)
	// or repository/helpers/search, keeping the annotation independent of folders.
	candidates:=[]string{parts[0], filepath.Join("helpers",parts[0])}
	for _,relative:=range candidates {
		candidate:=filepath.Join(repositoryDir,relative)
		files,err:=layout.GoFiles(candidate)
		if err!=nil{return "","","","",err}
		for _,f:=range files {
			pkg,err:=parser.ParseFile(token.NewFileSet(),f,nil,parser.PackageClauseOnly)
			if err!=nil{return "","","","",err}
			if pkg.Name.Name!=parts[0] {continue}
			prefix:="repository/"
            if strings.Contains(filepath.ToSlash(repositoryDir),"/internal/app/repository"){prefix="internal/app/repository/"}
			if filepath.Base(repositoryDir)!="repository"{prefix+=filepath.Base(repositoryDir)+"/"}
			return candidate,modulePath+"/services/"+service+"/"+prefix+filepath.ToSlash(relative),parts[0],parts[1],nil
		}
	}
	return "","","","",fmt.Errorf("@Searchable config package %q not found under %s or its helpers directory",parts[0],repositoryDir)
}
