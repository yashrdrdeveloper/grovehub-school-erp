// Package contract owns annotation metadata without depending on processors.
package contract

// Status identifies the implementation status of a directive.
type Status string

const (
	Existing    Status = "EXISTING"
	Gombok      Status = "GOMBOK"
	Proposed    Status = "PROPOSED"
	Optional    Status = "OPTIONAL"
	Discouraged Status = "DISCOURAGED"
	Pilot       Status = "PILOT"
)
