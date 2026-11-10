// Command check prevents unimplemented annotation contracts from being used
// as if they enforce security, transactions, messaging or validation behavior.
//
// A cataloged directive is not executable until an existing/pilot processor
// owns it. This command is a compile-time safety gate, not a runtime policy.
package main

import (
    "fmt"
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
    return nil
}

// CheckSource validates a single maintained Go file and is useful for
// annotation tooling and editor integrations without invoking the build.
func CheckSource(path string) error {return checkFile(path)}

