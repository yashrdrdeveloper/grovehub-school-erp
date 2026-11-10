// Command check prevents unimplemented annotation contracts from being used
// as if they enforce security, transactions, messaging or validation behavior.
//
// A cataloged directive is not executable until an existing/pilot processor
// owns it. This command is a compile-time safety gate, not a runtime policy.
package main

import (
    "fmt"
    "go/ast"
    goparser "go/parser"
    "go/token"
    "reflect"
    "strconv"
    "io/fs"
    "os"
    "path/filepath"
    "strings"

    "grovehub/libs/core/annotation"
    sourceparser "grovehub/libs/core/annotation/generation/parser"
)

func main(){
    if len(os.Args)!=2 {fmt.Fprintln(os.Stderr,"usage: check <backend-directory>");os.Exit(2)}
    if err:=checkServices(os.Args[1]);err!=nil {fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}

func checkServices(backend string)error {
    services:=filepath.Join(backend,"services")
    return filepath.WalkDir(services,func(path string,entry fs.DirEntry,walkErr error)error{
        if walkErr!=nil{return walkErr}
        if entry.IsDir(){
            // Generated code and vendored code are not maintained directive sources.
            if entry.Name()==".generated"||entry.Name()=="vendor" {return filepath.SkipDir}
            return nil
        }
        if !strings.HasSuffix(path,".go")||strings.HasSuffix(path,"_test.go"){return nil}
        return checkFile(path)
    })
}

func checkFile(path string) error {
    occurrences,err:=sourceparser.ScanFile(path)
    if err!=nil{return err}
    for _,found:=range occurrences {
        definition,exists:=annotation.Lookup(found.Directive.Name)
        if !exists {return fmt.Errorf("%s:%d: unknown annotation @%s",found.File,found.Line,found.Directive.Name)}
        allowed:=false
        matchedTarget:=false
        for _,context:=range definition.UseCases {
            if context.Target!=found.Target{continue}
            matchedTarget=true
            if context.Status==annotation.Existing||context.Status==annotation.Pilot {
                allowed=true
                break
            }
        }
        if !allowed {
            if matchedTarget {
                return fmt.Errorf("%s:%d: @%s is catalog metadata only, not an enforced processor; remove it until implemented",found.File,found.Line,found.Directive.Name)
            }
            return fmt.Errorf("%s:%d: @%s cannot annotate target %s",found.File,found.Line,found.Directive.Name,found.Target)
        }
    }
    return checkValidationTags(path)
}

// checkValidationTags rejects validation comments without executable field tags.
func checkValidationTags(path string) error {
    file,err:=goparser.ParseFile(token.NewFileSet(),path,nil,goparser.ParseComments)
    if err!=nil{return err}
    var validationError error
    ast.Inspect(file,func(node ast.Node) bool {
        field,ok:=node.(*ast.Field)
        if !ok||field.Doc==nil||validationError!=nil{return validationError==nil}
        tags:=""
        if field.Tag!=nil {
            raw,err:=strconv.Unquote(field.Tag.Value)
            if err!=nil {validationError=err;return false}
            tags=reflect.StructTag(raw).Get("validate")
        }
        for _,comment:=range field.Doc.List {
            directive,found,err:=annotation.ParseComment(comment.Text)
            if err!=nil {validationError=err;return false}
            if !found {continue}
            rule:=""
            switch directive.Name {
            case "Valid":rule="valid"
            case "NotBlank":rule="notblank"
            case "NotNull","NonNull":rule="nonnull"
            case "Size":
                parts:=strings.Fields(directive.Arguments)
                if len(parts)!=1||!strings.HasPrefix(parts[0],"max=") {
                    validationError=fmt.Errorf("%s: @Size requires max=N",path);return false
                }
                size,err:=strconv.Atoi(strings.TrimPrefix(parts[0],"max="))
                if err!=nil||size<0 {validationError=fmt.Errorf("%s: invalid @Size maximum",path);return false}
                rule=fmt.Sprintf("maxrunes=%d",size)
            default:continue
            }
            matched:=false
            for _,tag:=range strings.Split(tags,","){if tag==rule{matched=true;break}}
            if !matched {validationError=fmt.Errorf("%s: @%s must have validate:%s",path,directive.Name,rule);return false}
        }
        return true
    })
    return validationError
}

// CheckSource validates a single maintained Go file and is useful for
// annotation tooling and editor integrations without invoking the build.
func CheckSource(path string) error {return checkFile(path)}

