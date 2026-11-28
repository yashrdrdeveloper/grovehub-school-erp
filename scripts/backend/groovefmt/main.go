// Command groovefmt applies gofmt and GroveHub's one-line signature style.
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var newlineWhitespace = regexp.MustCompile(`[ \t]*\n[ \t]*`)
var trailingComma = regexp.MustCompile(`,\s*\)`)
var openingParen = regexp.MustCompile(`\( +`)

type edit struct {
	start, end int
	value      string
}

func main() {
	args := os.Args[1:]
	check := len(args) > 0 && args[0] == "-check"
	if check {
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: groovefmt [-check] <file-or-directory>...")
		os.Exit(2)
	}
	for _, path := range args {
		err := filepath.WalkDir(path, func(name string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				switch d.Name() {
				case ".git", "vendor", ".tools":
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(name, ".go") {
				return nil
			}
			if check {
				return checkFile(name)
			}
			return formatFile(name)
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func formatFile(path string) error {
	original, result, err := formattedFile(path)
	if err != nil {
		return err
	}
	if bytes.Equal(original, result) {
		return nil
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	return os.WriteFile(path, result, info.Mode().Perm())
}

// checkFile is a non-mutating lint of repository function and interface
// signatures. gofmt is available separately through make backend-format.
func checkFile(path string) error {
	source, err := os.ReadFile(path)
	if err != nil { return err }
	if bytes.Contains(source, []byte("// Code generated")) && bytes.Contains(source, []byte("DO NOT EDIT")) { return nil }
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, source, parser.AllErrors)
	if err != nil { return err }
	isOneLine := func(start, end token.Pos) bool {
		return fset.Position(start).Line == fset.Position(end).Line
	}
	for _, declaration := range file.Decls {
		switch decl := declaration.(type) {
		case *ast.FuncDecl:
			if !isOneLine(decl.Pos(), decl.Type.End()) {
				return fmt.Errorf("%s:%d: function signature must be on one line", path, fset.Position(decl.Pos()).Line)
			}
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				typ, ok := spec.(*ast.TypeSpec)
				if !ok { continue }
				if iface, ok := typ.Type.(*ast.InterfaceType); ok {
					for _, method := range iface.Methods.List {
						if len(method.Names) == 0 { continue }
						fn, ok := method.Type.(*ast.FuncType)
						if ok && !isOneLine(method.Pos(), fn.End()) {
							return fmt.Errorf("%s:%d: interface method signature must be on one line", path, fset.Position(method.Pos()).Line)
						}
					}
				}
			}
		}
	}
	return nil
}

func formattedFile(path string) ([]byte, []byte, error) {
	original, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	if bytes.Contains(original, []byte("// Code generated")) && bytes.Contains(original, []byte("DO NOT EDIT")) {
		return original, original, nil
	}
	formatted, err := format.Source(original)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: gofmt: %w", path, err)
	}
	result, err := singleLineSignatures(formatted)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", path, err)
	}
	if _, err := parser.ParseFile(token.NewFileSet(), path, result, parser.AllErrors); err != nil {
		return nil, nil, fmt.Errorf("%s: custom formatting: %w", path, err)
	}
	return original, result, nil
}

func singleLineSignatures(source []byte) ([]byte, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "input.go", source, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	var edits []edit
	add := func(start, end token.Pos) {
		a, b := fset.Position(start).Offset, fset.Position(end).Offset
		if a >= b {
			return
		}
		section := source[a:b]
		if !bytes.ContainsRune(section, '\n') {
			return
		}
		// Line comments and multiline literals must retain their original boundaries.
		for _, group := range file.Comments {
			for _, comment := range group.List {
				if comment.Pos() < end && comment.End() > start {
					return
				}
			}
		}
		if bytes.Contains(section, []byte("`")) {
			return
		}
		line := newlineWhitespace.ReplaceAllString(string(section), " ")
		line = trailingComma.ReplaceAllString(line, ")")
		line = openingParen.ReplaceAllString(line, "(")
		edits = append(edits, edit{a, b, line})
	}
	for _, declaration := range file.Decls {
		switch d := declaration.(type) {
		case *ast.FuncDecl:
			add(d.Pos(), d.Type.End())
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				t, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				switch typ := t.Type.(type) {
				case *ast.InterfaceType:
					for _, method := range typ.Methods.List {
						if len(method.Names) != 0 {
							if _, ok := method.Type.(*ast.FuncType); ok {
								add(method.Pos(), method.End())
							}
						}
					}
				case *ast.FuncType:
					add(t.Type.Pos(), t.Type.End())
				}
			}
		}
	}
	sort.Slice(edits, func(i, j int) bool { return edits[i].start > edits[j].start })
	for _, e := range edits {
		source = append(append(append([]byte{}, source[:e.start]...), e.value...), source[e.end:]...)
	}
	// gofmt itself enforces blank lines between top-level func/type/const declarations.
	return source, nil
}
