// Package annotation defines the Phase 1 metadata contracts for GrooveHub
// source directives. Except for the pre-existing @Entity/Gombok toolchain,
// catalog entries do not execute behavior, register services, or modify builds.
//
// Go does not have native Java-style annotations. This package describes
// directives written as comments, for example: // @Mapper.
package annotation

import (
    "errors"
    "fmt"
    "sort"
    "strings"
    "unicode"
)

// Status is the implementation status in the original backend master plan.
// Phase 1 defines contracts; Proposed/Optional entries are not executable.
type Status string

const (
    Existing   Status = "EXISTING"
    Gombok     Status = "GOMBOK"
    Proposed   Status = "PROPOSED"
    Optional   Status = "OPTIONAL"
    Discouraged Status = "DISCOURAGED"
    Pilot Status = "PILOT" // processor exists but full backend verification is pending
)

// Mode describes the planned mechanism rather than runtime behavior.
// Combinations, such as GEN/CHECK or REG/WRAP, are intentionally preserved.
type Mode string

const (
    Generate Mode = "GEN"
    Register Mode = "REG"
    Check    Mode = "CHECK"
    Wrap     Mode = "WRAP"
    Tag      Mode = "TAG"
)

// Target identifies a legal structural source location in Phase 1.
// Specific targets are refined per directive by Phase 2 processors.
type Target string

const (
    TypeTarget      Target = "type"
    InterfaceTarget Target = "interface"
    MethodTarget    Target = "method"
    FunctionTarget  Target = "function"
    FieldTarget     Target = "field"
    PackageTarget   Target = "package"
    FileTarget      Target = "file"
    ParameterTarget Target = "parameter"
)

// UseCase is one catalog row. A directive can have multiple UseCases
// (for example @Retry appears in both external clients and middleware).
type UseCase struct {
    Name        string
    Syntax      string
    Category    string
    Purpose     string
    Mode        string
    Status      Status
    Roadmap     string
    Target      Target
}

// Definition groups every supported context for one directive.
// All returned definitions are owned by the caller (safe to modify).
type Definition struct {
    Name     string
    UseCases []UseCase
}

// Pattern is a GoF design-pattern annotation decision. A pattern's presence
// in the catalog does not imply that automatic pattern implementation is safe.
type Pattern struct {
    Name     string
    Family   string
    Pattern  string
    Purpose  string
    Decision string
}

// Directive is a parsed comment directive. Arguments stay uninterpreted
// until a dedicated, validated Phase 2 processor is implemented.
type Directive struct {
    Name      string
    Arguments string
}

var (
    ErrUnknownDirective = errors.New("unknown annotation directive")
    ErrInvalidDirective = errors.New("invalid annotation directive")
)

var known = buildDefinitions()

func buildDefinitions() map[string]Definition {
    registry := make(map[string]Definition)
    for _, use := range inventory {
        if use.Name == "" || use.Syntax == "" || use.Category == "" || use.Purpose == "" || use.Target == "" {
            panic(fmt.Sprintf("incomplete annotation definition: %+v", use))
        }
        if !strings.HasPrefix(use.Syntax, "@"+use.Name) {
            panic("annotation syntax/name mismatch: " + use.Name)
        }
        d := registry[use.Name]
        d.Name = use.Name
        d.UseCases = append(d.UseCases, use)
        registry[use.Name] = d
    }
    // Every GoF candidate annotation also has a discoverable declaration,
    // including patterns intentionally not recommended for generation.
    for _, pattern := range patternInventory {
        if _, exists := registry[pattern.Name]; exists { continue }
        registry[pattern.Name] = Definition{
            Name: pattern.Name,
            UseCases: []UseCase{{
                Name:pattern.Name, Syntax:"@"+pattern.Name,
                Category:"GoF design patterns", Purpose:pattern.Purpose,
                Mode:"CHECK", Status:Optional, Roadmap:"P8", Target:TypeTarget,
            }},
        }
    }
    return registry
}

// All lists every unique annotation in deterministic lexicographic order.
func All() []Definition {
    names := make([]string, 0, len(known))
    for name := range known { names = append(names,name) }
    sort.Strings(names)
    result := make([]Definition,0,len(names))
    for _, name := range names {
        entry := known[name]
        entry.UseCases = append([]UseCase(nil),entry.UseCases...)
        result = append(result,entry)
    }
    return result
}

// Lookup returns every context for a single annotation; names exclude '@'.
func Lookup(name string) (Definition,bool) {
    name = strings.TrimPrefix(name,"@")
    entry, ok := known[name]
    if !ok { return Definition{},false }
    entry.UseCases = append([]UseCase(nil),entry.UseCases...)
    return entry,true
}

// Uses lists every row, preserving multi-context annotations such as
// @Timeout, @Sensitive, @Retry, and @SoftDelete.
func Uses() []UseCase { return append([]UseCase(nil), inventory...) }

// Patterns lists all 23 GoF pattern decisions, including discouraged ones.
func Patterns() []Pattern { return append([]Pattern(nil),patternInventory...) }

// ParseComment reads a single Go line comment such as:
//    // @Entity foundation.persons
// or // @NonNull identity.InvalidRegistration
// This is metadata parsing, not executable validation or code generation.
// Non-annotation comments return (Directive{}, false, nil).
func ParseComment(line string) (Directive,bool,error) {
    line = strings.TrimSpace(line)
    if !strings.HasPrefix(line,"//") { return Directive{},false,nil }
    line = strings.TrimSpace(strings.TrimPrefix(line,"//"))
    if !strings.HasPrefix(line,"@") { return Directive{},false,nil }
    line = strings.TrimPrefix(line,"@")
    index := strings.IndexFunc(line,unicode.IsSpace)
    name := line
    args := ""
    if index >= 0 {
        name = line[:index]
        args = strings.TrimSpace(line[index:])
    }
    if name=="" {return Directive{},true,fmt.Errorf("%w: name is empty",ErrInvalidDirective)}
    for _, ch := range name {
        if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
            return Directive{},true,fmt.Errorf("%w: invalid name %q",ErrInvalidDirective,name)
        }
    }
    if _,ok := Lookup(name); !ok {
        return Directive{},true,fmt.Errorf("%w: @%s",ErrUnknownDirective,name)
    }
    d:=Directive{Name:name,Arguments:args}
    if err:=ValidateDirective(d);err!=nil {return Directive{},true,err}
    return d,true,nil
}

// ValidateDirective checks the catalog and Phase 1's universally safe syntax
// constraints. The domain-specific meaning of parameters belongs to
// individual processors and will not be guessed by this package.
func ValidateDirective(d Directive) error {
    if _,ok:=Lookup(d.Name);!ok {
        return fmt.Errorf("%w: @%s",ErrUnknownDirective,d.Name)
    }
    if strings.ContainsAny(d.Arguments,"\r\n") {
        return fmt.Errorf("%w: multiline arguments are not supported",ErrInvalidDirective)
    }
    if d.Name=="Entity" {
        parts:=strings.Split(d.Arguments,".")
        if len(parts)!=2 || !identifier(parts[0]) || !identifier(parts[1]) {
            return fmt.Errorf("%w: @Entity requires schema.table",ErrInvalidDirective)
        }
    }
    return nil
}

func identifier(value string) bool {
    if value=="" {return false}
    for i,ch:= range value {
        if i==0 {
            if !(ch>='a'&&ch<='z') && ch!='_' {return false}
        } else if !(ch>='a'&&ch<='z') && !(ch>='0'&&ch<='9') && ch!='_' {
            return false
        }
    }
    return true
}
