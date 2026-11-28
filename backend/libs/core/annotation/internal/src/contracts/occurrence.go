package contract

// Occurrence identifies a directive and its source location.
type Occurrence struct {
	Directive Directive
	Target    Target
	File      string
	Line      int
	Column    int
}
