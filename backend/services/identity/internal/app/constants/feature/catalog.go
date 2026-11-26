// Package feature owns the registered feature code vocabulary.
package feature

type Code string

const (
	AccountManagement Code = "IDENTITY_ACCOUNT_MANAGEMENT"
	RoleManagement    Code = "IDENTITY_ROLE_MANAGEMENT"
)

func Known(code Code) bool {
	switch code {
	case AccountManagement, RoleManagement:
		return true
	default:
		return false
	}
}
