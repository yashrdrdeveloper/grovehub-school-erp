// Package feature owns the registered feature code vocabulary.
package feature

type Code string

const (
	AccountManagement Code = "IDENTITY_ACCOUNT_MANAGEMENT"
	RoleManagement    Code = "IDENTITY_ROLE_MANAGEMENT"
)

type Definition struct {
	Code Code
	Name string
}

func Catalog() []Definition {
	return []Definition{
		{Code: AccountManagement, Name: "Account management"},
		{Code: RoleManagement, Name: "Role management"},
	}
}

func Find(code Code) (Definition, bool) {
	for _, item := range Catalog() {
		if item.Code == code {
			return item, true
		}
	}
	return Definition{}, false
}

func Known(code Code) bool {
	_, known := Find(code)
	return known
}
