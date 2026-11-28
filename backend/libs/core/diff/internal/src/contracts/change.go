package contracts

// Change describes the old and new JSON values at a field path.
type Change struct {
	Field    string `json:"field"`
	OldValue any    `json:"oldValue"`
	NewValue any    `json:"newValue"`
}
