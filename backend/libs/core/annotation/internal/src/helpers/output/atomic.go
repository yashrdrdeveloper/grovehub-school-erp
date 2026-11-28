package output

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// TempFile is the file boundary for atomic replacement. Using an interface
// allows failures during writing and closing to be tested without disk faults.
type TempFile interface {
	io.WriteCloser
	Name() string
}

type AtomicWriter struct {
	CreateTemp func(dir, pattern string) (TempFile,error)
	Rename func(oldPath,newPath string) error
	Remove func(path string) error
}

func WriteAtomic(path, pattern string, contents []byte) error {
	writer:=AtomicWriter{
		CreateTemp:func(dir,pattern string)(TempFile,error){return os.CreateTemp(dir,pattern)},
		Rename:os.Rename,
		Remove:os.Remove,
	}
	return writer.Write(path,pattern,contents)
}

// Write only replaces the destination after a successful write and close.
// Failed attempts remove their temporary file and preserve the old destination.
func (writer AtomicWriter) Write(path, pattern string, contents []byte) error {
	temp,err:=writer.CreateTemp(filepath.Dir(path),pattern)
	if err!=nil{return err}
	defer writer.Remove(temp.Name())
	_,writeErr:=temp.Write(contents)
	closeErr:=temp.Close()
	if writeErr!=nil{return writeErr}
	if closeErr!=nil{return closeErr}
	return writer.Rename(temp.Name(),path)
}

// WriteJSON writes an indented deterministic JSON document atomically.
// The old document remains untouched if serialization or persistence fails.
func WriteJSON(path,pattern string,value any)error{
	if err:=os.MkdirAll(filepath.Dir(path),0755);err!=nil{return err}
	data,err:=json.MarshalIndent(value,"","  ")
	if err!=nil{return err}
	data=append(data,'\n')
	old,err:=os.ReadFile(path)
	if err==nil&&bytes.Equal(old,data){return nil}
	if err!=nil&&!os.IsNotExist(err){return err}
	return WriteAtomic(path,pattern,data)
}
