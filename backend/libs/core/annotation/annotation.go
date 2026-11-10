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

    "grovehub/libs/core/annotation/contract"
    "grovehub/libs/core/annotation/patterns"
)

// Re-export the contract types to preserve the annotation package's API.
// The category packages depend only on contract, never on this registry.
type Status = contract.Status
const (
    Existing = contract.Existing
    Gombok = contract.Gombok
    Proposed = contract.Proposed
    Optional = contract.Optional
    Discouraged = contract.Discouraged
    Pilot = contract.Pilot
)
type Mode = contract.Mode
const (
    Generate = contract.Generate
    Register = contract.Register
    Check = contract.Check
    Wrap = contract.Wrap
    Tag = contract.Tag
)
type Target = contract.Target
const (
    TypeTarget = contract.TypeTarget
    InterfaceTarget = contract.InterfaceTarget
    MethodTarget = contract.MethodTarget
    FunctionTarget = contract.FunctionTarget
    FieldTarget = contract.FieldTarget
    PackageTarget = contract.PackageTarget
    FileTarget = contract.FileTarget
    ParameterTarget = contract.ParameterTarget
)
type UseCase = contract.UseCase

// Definition groups every supported context for one directive.
// All returned definitions are owned by the caller (safe to modify).
type Definition struct {
    Name     string
    UseCases []UseCase
}

// A pattern is a GoF architectural decision, not a generated service.
type Pattern = contract.Pattern

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

var patternInventory = patterns.Entries

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
        implementation := Optional
        mode := "CHECK"
        target := TypeTarget
        if pattern.Name == "Proxy" {
            implementation = Pilot
            mode = "GEN"
            target = MethodTarget
        }
        registry[pattern.Name] = Definition{
            Name: pattern.Name,
            UseCases: []UseCase{{
                Name:pattern.Name, Syntax:"@"+pattern.Name,
                Category:"GoF design patterns", Purpose:pattern.Purpose,
                Mode:mode, Status:implementation, Roadmap:"P8", Target:target,
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
        args:=strings.Fields(d.Arguments)
        if len(args)<1 || len(args)>2 {
            return fmt.Errorf("%w: @Entity requires schema.table [mode=table]",ErrInvalidDirective)
        }
        parts:=strings.Split(args[0],".")
        if len(parts)!=2 || !identifier(parts[0]) || !identifier(parts[1]) {
            return fmt.Errorf("%w: @Entity requires schema.table",ErrInvalidDirective)
        }
        if len(args)==2 && args[1]!="mode=table" {
            return fmt.Errorf("%w: unsupported @Entity option %q",ErrInvalidDirective,args[1])
        }
    }
    switch d.Name {
    case "Repository":
        parts, err := directiveOptions(d.Arguments, []string{"entity"})
        if err != nil || !tokenIdentifier(parts["entity"]) {
            return fmt.Errorf("%w: @Repository requires entity=Type", ErrInvalidDirective)
        }
    case "Searchable":
        parts, err := directiveOptions(d.Arguments, []string{"impl", "entity", "config", "alias"})
        if err != nil || !tokenIdentifier(parts["impl"]) || !tokenIdentifier(parts["entity"]) ||
            !tokenIdentifier(parts["config"]) || !tokenIdentifier(parts["alias"]) {
            return fmt.Errorf("%w: @Searchable requires impl=Type entity=Type config=identifier alias=identifier", ErrInvalidDirective)
        }
    case "Filterable", "Sortable":
        parts, err := directiveOptions(d.Arguments, []string{"fields"})
        if err != nil || !validFieldList(parts["fields"]) {
            return fmt.Errorf("%w: @%s requires fields=comma,separated,names", ErrInvalidDirective, d.Name)
        }
    case "Cursor":
        parts, err := directiveOptions(d.Arguments, []string{"field"})
        if err != nil || !validFieldName(parts["field"]) {
            return fmt.Errorf("%w: @Cursor requires field=fieldName", ErrInvalidDirective)
        }
    }
    return nil
}

func directiveOptions(raw string, required []string) (map[string]string,error) {
    result := make(map[string]string,len(required))
    for _, arg := range strings.Fields(raw) {
        pieces := strings.SplitN(arg,"=",2)
        if len(pieces)!=2 || pieces[1]=="" {return nil,ErrInvalidDirective}
        allowed := false
        for _, key := range required {if pieces[0]==key {allowed=true;break}}
        if !allowed || result[pieces[0]]!="" {return nil,ErrInvalidDirective}
        result[pieces[0]]=pieces[1]
    }
    for _, key := range required {if result[key]=="" {return nil,ErrInvalidDirective}}
    return result,nil
}
func tokenIdentifier(s string) bool {
    if s=="" {return false}
    for i,r := range s {
        if i==0 {if !(r>='a'&&r<='z'||r>='A'&&r<='Z'||r=='_'){return false}}
        if !(r>='a'&&r<='z'||r>='A'&&r<='Z'||r>='0'&&r<='9'||r=='_'){return false}
    }
    return true
}
func validFieldName(s string) bool {
    if s=="" {return false}
    for _,part := range strings.Split(s,".") {if !tokenIdentifier(part) {return false}}
    return true
}
func validFieldList(s string) bool {
    if s=="" {return false}
    seen:=map[string]bool{}
    for _,field:=range strings.Split(s,",") {
        if !validFieldName(field)||seen[field] {return false}
        seen[field]=true
    }
    return true
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
