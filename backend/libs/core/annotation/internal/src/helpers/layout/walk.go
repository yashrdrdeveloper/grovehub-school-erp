package layout

import (
	"io/fs"
	"os"
	"path/filepath"
)

// WalkOwned walks an existing module-relative directory without following
// directory symlinks. The visitor receives the relative path from the owning
// module and can return fs.SkipDir. Missing optional roots are not errors.
func WalkOwned(owner, relative string, visit func(path, relative string, entry fs.DirEntry)error)error{
	root,err:=ResolveRelative(owner,relative)
	if err!=nil{return err}
	info,err:=os.Lstat(root)
	if os.IsNotExist(err){return nil}
	if err!=nil{return err}
	// Report inaccessible and cyclic root symlinks as errors, but do not
	// traverse a valid directory symlink or a missing symlink target.
	if info.Mode()&os.ModeSymlink!=0{
		_,err:=os.Stat(root)
		if os.IsNotExist(err){return nil}
		return err
	}
	if !info.IsDir(){return nil}
	return fs.WalkDir(os.DirFS(owner),filepath.ToSlash(filepath.Clean(relative)),func(relative string,entry fs.DirEntry,walkErr error)error{
		if walkErr!=nil{return walkErr}
		local:=filepath.FromSlash(relative)
		return visit(filepath.Join(owner,local),local,entry)
	})
}

// WalkOwnedDirectories visits only maintained directory packages, optionally
// skipping directory names such as vendor and .generated.
func WalkOwnedDirectories(owner,relative string, skip map[string]bool,visit func(path,relative string)error)error{
	return WalkOwned(owner,relative,func(path,relative string,entry fs.DirEntry)error{
		if !entry.IsDir(){return nil}
		if skip[entry.Name()]{return fs.SkipDir}
		return visit(path,relative)
	})
}
