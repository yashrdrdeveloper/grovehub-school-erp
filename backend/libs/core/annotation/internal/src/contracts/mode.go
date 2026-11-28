package contract

// Mode describes a directive's processing mechanism.
type Mode string

const (
	Generate Mode = "GEN"
	Register Mode = "REG"
	Check    Mode = "CHECK"
	Wrap     Mode = "WRAP"
	Tag      Mode = "TAG"
)
