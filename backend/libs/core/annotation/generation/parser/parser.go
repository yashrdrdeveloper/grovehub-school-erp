// Package parser discovers Phase 1 annotation comments in Go source files.
// Discovery never rewrites files, generates implementations, or executes policy.
// Future per-capability processors can consume these typed occurrences.
package parser

import (
    "fmt"
    "go/ast"
    goparser "go/parser"
    "go/token"

    "grovehub/libs/core/annotation"
)

// Occurrence identifies a known annotation and its source declaration.
// Target is advisory until a directive-specific processor is available.
type Occurrence struct {
    Directive annotation.Directive
    Target    annotation.Target
    File      string
    Line      int
    Column    int
}

// ScanFile parses a Go file using the standard AST and returns recognized
// directives in source order. It fails closed on unknown or invalid directives.
func ScanFile(path string) ([]Occurrence,error) {
    set:=token.NewFileSet()
    file,err:=goparser.ParseFile(set,path,nil,goparser.ParseComments|goparser.AllErrors)
    if err!=nil{return nil,err}

    targets:=map[*ast.CommentGroup]annotation.Target{}
    ast.Inspect(file,func(n ast.Node) bool {
        switch decl:=n.(type) {
        case *ast.File:
            if decl.Doc!=nil {targets[decl.Doc]=annotation.PackageTarget}
        case *ast.FuncDecl:
            if decl.Doc!=nil {
                t:=annotation.FunctionTarget
                if decl.Recv!=nil {t=annotation.MethodTarget}
                targets[decl.Doc]=t
            }
        case *ast.GenDecl:
            if decl.Doc!=nil {
                target:=annotation.FileTarget
                if decl.Tok==token.TYPE {
                    target=annotation.TypeTarget
                    if len(decl.Specs)==1 {
                        if spec,ok:=decl.Specs[0].(*ast.TypeSpec);ok {
                            if _,ok:=spec.Type.(*ast.InterfaceType);ok {target=annotation.InterfaceTarget}
                        }
                    }
                }
                targets[decl.Doc]=target
            }
        case *ast.TypeSpec:
            if decl.Doc!=nil {
                target:=annotation.TypeTarget
                if _,ok:=decl.Type.(*ast.InterfaceType);ok {target=annotation.InterfaceTarget}
                targets[decl.Doc]=target
            }
        case *ast.Field:
            if decl.Doc!=nil {targets[decl.Doc]=annotation.FieldTarget}
            if decl.Comment!=nil {targets[decl.Comment]=annotation.FieldTarget}
        }
        return true
    })

    var result []Occurrence
    for _,group:=range file.Comments {
        target,ok:=targets[group]
        if !ok {target=annotation.FileTarget}
        for _,comment:=range group.List {
            d,found,err:=annotation.ParseComment(comment.Text)
            if err!=nil {
                pos:=set.Position(comment.Pos())
                return nil,fmt.Errorf("%s:%d:%d: %w",pos.Filename,pos.Line,pos.Column,err)
            }
            if !found {continue}
            pos:=set.Position(comment.Pos())
            result=append(result,Occurrence{Directive:d,Target:target,File:pos.Filename,Line:pos.Line,Column:pos.Column})
        }
    }
    return result,nil
}
