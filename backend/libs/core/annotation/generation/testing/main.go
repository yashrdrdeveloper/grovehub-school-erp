// Command testing generates Ginkgo suite registration and test-owned GoMock
// implementations. No generated testing code is imported by production packages.
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
    "os/exec"
    "path/filepath"
    "sort"
    "strings"
    "unicode"

    "grovehub/libs/core/annotation"
)

type testSuite struct {
    name string
    pkg string
    specs []string
    beforeAll string
    beforeEach []string
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: testing <backend-directory>")
        os.Exit(2)
    }
    if err := generateAll(os.Args[1]); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func generateAll(backend string) error {
    backend, err := filepath.Abs(backend)
    if err != nil { return err }
    roots := []struct{ owner, tests string }{
        {filepath.Join(backend,"libs","core"), filepath.Join(backend,"libs","core","tests")},
    }
    services, err := os.ReadDir(filepath.Join(backend,"services"))
    if err != nil { return err }
    for _, service := range services {
        if !service.IsDir() { continue }
        owner := filepath.Join(backend,"services",service.Name())
        roots = append(roots, struct{owner,tests string}{owner,filepath.Join(owner,"internal","tests")})
    }
    for _, root := range roots {
        if _, err := os.Stat(root.tests); os.IsNotExist(err) { continue } else if err != nil { return err }
        if err := filepath.WalkDir(root.tests,func(dir string,entry fs.DirEntry,walkErr error) error {
            if walkErr != nil { return walkErr }
            if !entry.IsDir() { return nil }
            if entry.Name() == ".generated" || entry.Name() == "vendor" { return filepath.SkipDir }
            return generateDir(backend,root.owner,dir)
        }); err != nil { return err }
    }
    return nil
}

func generateDir(backend, owner, dir string) error {
    paths, err := filepath.Glob(filepath.Join(dir,"*.go"))
    if err != nil { return err }
    sort.Strings(paths)
    var suite testSuite
    var mockSource string
    for _, path := range paths {
        if strings.HasPrefix(filepath.Base(path),"zz_") { continue }
        f, err := parser.ParseFile(token.NewFileSet(),path,nil,parser.ParseComments)
        if err != nil { return fmt.Errorf("%s: %w",path,err) }
        isTest := strings.HasSuffix(path,"_test.go")
        for _, declaration := range f.Decls {
            if group,ok := declaration.(*ast.GenDecl); ok && group.Tok == token.TYPE {
                for _, item := range group.Specs {
                    spec,ok := item.(*ast.TypeSpec)
                    if !ok { continue }
                    doc := spec.Doc
                    if doc == nil { doc = group.Doc }
                    args,marked,err := directive(doc,"Mock")
                    if err != nil { return fmt.Errorf("%s: %w",path,err) }
                    if !marked { continue }
                    if args != "" { return fmt.Errorf("%s: @Mock does not accept arguments",path) }
                    if _,ok := spec.Type.(*ast.InterfaceType); !ok { return fmt.Errorf("%s: @Mock requires interface %s",path,spec.Name.Name) }
                    if isTest { return fmt.Errorf("%s: @Mock interface must be in a normal Go file in internal/tests/mocks",path) }
                    if !strings.Contains(filepath.ToSlash(dir),"/tests/mocks/") { return fmt.Errorf("%s: @Mock is limited to test-owned mocks",path) }
                    if mockSource != "" && mockSource != path { return fmt.Errorf("%s: at most one @Mock source file per directory",dir) }
                    mockSource = path
                }
            }
            fn,ok := declaration.(*ast.FuncDecl)
            if !ok { continue }
            for _, kind := range []string{"Testing","BeforeAll","BeforeEach"} {
                args,marked,err := directive(fn.Doc,kind)
                if err != nil { return fmt.Errorf("%s: %w",path,err) }
                if !marked { continue }
                if !isTest { return fmt.Errorf("%s: @%s must decorate a *_test.go function",path,kind) }
                if fn.Recv != nil || fn.Type.TypeParams != nil ||
                    (fn.Type.Params != nil && len(fn.Type.Params.List)>0) ||
                    (fn.Type.Results != nil && len(fn.Type.Results.List)>0) {
                    return fmt.Errorf("%s: @%s function %s must take no parameters and return nothing",path,kind,fn.Name.Name)
                }
                if suite.pkg != "" && suite.pkg != f.Name.Name { return fmt.Errorf("%s: inconsistent test package",dir) }
                suite.pkg = f.Name.Name
                switch kind {
                case "Testing":
                    if !strings.HasPrefix(args,"suite=") || strings.ContainsAny(args," \t\r\n") || len(args)<=len("suite=") {
                        return fmt.Errorf("%s: @Testing requires exactly suite=Name",path)
                    }
                    if suite.name != "" && suite.name != strings.TrimPrefix(args,"suite=") {
                        return fmt.Errorf("%s: multiple @Testing suite names in package",dir)
                    }
                    suite.name = strings.TrimPrefix(args,"suite=")
                    suite.specs = append(suite.specs,fn.Name.Name)
                case "BeforeAll":
                    if args != "" { return fmt.Errorf("%s: @BeforeAll does not accept arguments",path) }
                    if suite.beforeAll != "" { return fmt.Errorf("%s: multiple @BeforeAll hooks",dir) }
                    suite.beforeAll = fn.Name.Name
                case "BeforeEach":
                    if args != "" { return fmt.Errorf("%s: @BeforeEach does not accept arguments",path) }
                    suite.beforeEach = append(suite.beforeEach,fn.Name.Name)
                }
            }
        }
    }
    rel,err := filepath.Rel(owner,dir)
    if err != nil { return err }
    outDir := filepath.Join(owner,".generated",rel)
    mockOutput := filepath.Join(outDir,"zz_mocks.gen.go")
    if mockSource != "" {
        f,err := parser.ParseFile(token.NewFileSet(),mockSource,nil,parser.PackageClauseOnly)
        if err != nil { return err }
        if err := os.MkdirAll(outDir,0755); err != nil { return err }
        cmd := exec.Command("go","run","go.uber.org/mock/mockgen","-source="+mockSource,"-destination="+mockOutput,"-package="+f.Name.Name)
        cmd.Dir = backend
        if output,err := cmd.CombinedOutput();err != nil { return fmt.Errorf("%s: generate mocks: %w: %s",mockSource,err,output) }
    } else if err := os.Remove(mockOutput); err != nil && !os.IsNotExist(err) { return err }

    suiteOutput := filepath.Join(outDir,"zz_suite_test.go")
    if suite.name != "" {
        for _,ch := range suite.name { if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) { return fmt.Errorf("%s: suite name must contain only letters and digits",dir) } }
        if unicode.IsDigit([]rune(suite.name)[0]) { return fmt.Errorf("%s: suite name must start with a letter",dir) }
        return writeSuite(suiteOutput,suite)
    }
    if suite.beforeAll != "" || len(suite.beforeEach)>0 { return fmt.Errorf("%s: testing hook requires an @Testing suite",dir) }
    if err := os.Remove(suiteOutput);err != nil && !os.IsNotExist(err) { return err }
    return nil
}

func directive(group *ast.CommentGroup,name string) (string,bool,error) {
    if group==nil { return "",false,nil }
    var result string
    found := false
    for _,comment := range group.List {
        parsed,recognized,err := annotation.ParseComment(comment.Text)
        if err != nil { return "",false,err }
        if !recognized || parsed.Name!=name { continue }
        if found { return "",false,fmt.Errorf("duplicate @%s",name) }
        result,found=parsed.Arguments,true
    }
    return result,found,nil
}

func writeSuite(output string,suite testSuite) error {
    var source strings.Builder
    fmt.Fprintln(&source,"// Code generated by GrooveHub @Testing. DO NOT EDIT.")
    fmt.Fprintf(&source,"package %s\n\n",suite.pkg)
    source.WriteString("import (\n\"testing\"\n. \"github.com/onsi/ginkgo/v2\"\n. \"github.com/onsi/gomega\"\n)\n")
    fmt.Fprintf(&source,"func Test%s(t *testing.T) {\n RegisterFailHandler(Fail)\n RunSpecs(t,%q)\n}\n",suite.name,suite.name)
    if suite.beforeAll!="" { fmt.Fprintf(&source,"var _ = BeforeSuite(%s)\n",suite.beforeAll) }
    fmt.Fprintf(&source,"var _ = Describe(%q,func(){\n",suite.name)
    for _,before := range suite.beforeEach { fmt.Fprintf(&source,"BeforeEach(%s)\n",before) }
    for _,spec := range suite.specs { fmt.Fprintf(&source,"%s()\n",spec) }
    source.WriteString("})\n")
    formatted,err := format.Source([]byte(source.String()))
    if err != nil { return err }
    old,err := os.ReadFile(output)
    if err==nil && bytes.Equal(old,formatted) { return nil }
    if err!=nil && !os.IsNotExist(err) { return err }
    if err := os.MkdirAll(filepath.Dir(output),0755); err != nil { return err }
    return os.WriteFile(output,formatted,0644)
}
