package contracts

// CursorValueValidator checks the decoded sort value or cursor ID using
// the error supplied by the owning service.
type CursorValueValidator func(value string, invalid error) error

// ValueNormalizer is resource-owned. Generic filter validation invokes it
// after checking fields, operators and value counts.
type ValueNormalizer func(operator Operator, value string) (string, error)

type FieldSpec struct {
    Operators map[Operator]struct{}
    Normalizer ValueNormalizer
}
