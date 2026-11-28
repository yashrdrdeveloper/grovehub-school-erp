package layout

import (
	"os"
	"path/filepath"
	"strings"
)

// GoFiles returns lexically ordered Go source paths from a literal directory.
// Unlike filepath.Glob, it neither interprets metacharacters in directory
// names nor silently hides read failures. Optional absent roots are empty.
func GoFiles(dir string)([]string,error){
	entries,err:=os.ReadDir(dir)
	if os.IsNotExist(err){return nil,nil}
	if err!=nil{return nil,err}
	var files []string
	for _,entry:=range entries{
		if strings.HasSuffix(entry.Name(),".go"){
			files=append(files,filepath.Join(dir,entry.Name()))
		}
	}
	return files,nil
}
