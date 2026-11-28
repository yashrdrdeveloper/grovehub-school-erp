// Command expand converts @Entity declarations into Gombok annotations in staging.
package main

import (
 "fmt"
 "go/ast"
 "go/parser"
 "go/token"
 "os"
 "path/filepath"
 "strings"
)

type methodKinds struct { getter, setter bool }

func main() {
 if len(os.Args)!=2 { fmt.Fprintln(os.Stderr,"usage: expand <staging-entity-directory>"); os.Exit(2) }
 if err:=expand(os.Args[1]);err!=nil {fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}

func expand(dir string) error {
 files,err:=filepath.Glob(filepath.Join(dir,"*.go"));if err!=nil{return err}
 fset:=token.NewFileSet()
 parsed:=make(map[string]*ast.File)
 existing:=make(map[string]methodKinds)
 for _,path:=range files {
  if strings.HasSuffix(path,"_test.go") || strings.HasSuffix(path,"_gombok.go") {continue}
  file,err:=parser.ParseFile(fset,path,nil,parser.ParseComments);if err!=nil{return err}
  parsed[path]=file
  for _,decl:=range file.Decls {
   fn,ok:=decl.(*ast.FuncDecl);if !ok||fn.Recv==nil||len(fn.Recv.List)!=1 {continue}
   recv:=fn.Recv.List[0].Type
   if ptr,ok:=recv.(*ast.StarExpr);ok {recv=ptr.X}
   name,ok:=recv.(*ast.Ident);if !ok {continue}
   kind:=existing[name.Name]
   if strings.HasPrefix(fn.Name.Name,"Get") {kind.getter=true}
   if strings.HasPrefix(fn.Name.Name,"Set") {kind.setter=true}
   existing[name.Name]=kind
  }
 }
 for path,file:=range parsed {
  data,err:=os.ReadFile(path);if err!=nil{return err}
  // Offset-based replacement preserves original source and Gombok's doc attachment.
  type edit struct{ start,end int; text string }
  var edits []edit
  for _,decl:=range file.Decls {
   group,ok:=decl.(*ast.GenDecl);if !ok||group.Tok!=token.TYPE {continue}
   for _,spec:=range group.Specs {
    typ,ok:=spec.(*ast.TypeSpec);if !ok {continue}
    if _,ok:=typ.Type.(*ast.StructType);!ok {continue}
    doc:=typ.Doc;if doc==nil {doc=group.Doc};if doc==nil {continue}
    for _,comment:=range doc.List {
     annotation:=strings.TrimSpace(strings.TrimPrefix(comment.Text,"//"))
     if !strings.HasPrefix(annotation,"@Entity ") {continue}
     table:=strings.TrimSpace(strings.TrimPrefix(annotation,"@Entity "))
     if table=="" || strings.ContainsAny(table," \t\r\n") {return fmt.Errorf("%s: @Entity requires a single schema.table name",typ.Name.Name)}
     kinds:=existing[typ.Name.Name]
     replacement:="// @Builder"
     if !kinds.getter {replacement+="\n// @Getter"}
     if !kinds.setter {replacement+="\n// @Setter"}
     start:=fset.Position(comment.Pos()).Offset
     end:=fset.Position(comment.End()).Offset
     edits=append(edits,edit{start,end,replacement})
    }
   }
  }
  for i:=len(edits)-1;i>=0;i-- {
   e:=edits[i]
   data=append(append(append([]byte{},data[:e.start]...),[]byte(e.text)...),data[e.end:]...)
  }
  if len(edits)>0 {
   if err:=os.WriteFile(path,data,0644);err!=nil{return err}
  }
 }
 return nil
}
