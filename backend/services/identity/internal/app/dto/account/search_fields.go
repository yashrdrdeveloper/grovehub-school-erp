package account

// Canonical API search keys. Repository SQL expressions and validation specs
// must use these names; DTO JSON-tag parity is covered by tests.
const (
	FieldAccountID = "accountId"
	FieldPersonGID = "personGid"
	FieldUsername = "username"
	FieldEmail = "email"
	FieldStatus = "status"
	FieldCreatedAt = "createdAt"
	FieldUpdatedAt = "updatedAt"
)
