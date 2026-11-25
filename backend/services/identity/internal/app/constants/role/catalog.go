package role

import accountdto "grovehub/services/identity/internal/app/dto/account"

// Catalog is the starter set of operational school roles. It does not grant
// gateway administrator authority or define permissions. Every assignment is
// scoped to an account belonging to the authenticated administrator's tenant.
func Catalog() []accountdto.Role {
	return []accountdto.Role{
		{ID: "SCHOOL_ADMINISTRATION", Name: "School administration"},
		{ID: "TEACHER", Name: "Teacher"},
		{ID: "COORDINATOR", Name: "Coordinator"},
		{ID: "ACCOUNTANT", Name: "Accountant"},
		{ID: "LIBRARIAN", Name: "Librarian"},
		{ID: "RECEPTIONIST", Name: "Receptionist"},
	}
}

func IsAssignable(id string) bool {
	for _, item := range Catalog() {
		if item.ID == id {
			return true
		}
	}
	return false
}
