// Package contract contains shared annotation metadata types without any
// dependency on the registry, generators, services or runtime policies.
package contract

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

// Pattern is a GoF design-pattern annotation decision. A pattern's presence
// in the catalog does not imply that automatic pattern implementation is safe.
type Pattern struct {
    Name     string
    Family   string
    Pattern  string
    Purpose  string
    Decision string
}
