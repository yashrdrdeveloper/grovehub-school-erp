package validation

const (
	ProjectCodePattern = "^[a-z][a-z0-9_-]{1,63}$"
	MaxNameBytes = 128
	MaxReferenceBytes = 128
	InitialVersion int64 = 1
	InitialGeneration int64 = 1
)
