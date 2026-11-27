// Package persona owns the fixed school job categories, not authorization grants.
package persona

type Code string

const (
	SchoolAdministration Code = "SCHOOL_ADMINISTRATION"
	Teacher              Code = "TEACHER"
	Coordinator          Code = "COORDINATOR"
	Accountant           Code = "ACCOUNTANT"
	Librarian            Code = "LIBRARIAN"
	Receptionist         Code = "RECEPTIONIST"
)

type Definition struct {
	Code Code
	Name string
}

func Catalog() []Definition {
	return []Definition{
		{Code: SchoolAdministration, Name: "School administration"},
		{Code: Teacher, Name: "Teacher"},
		{Code: Coordinator, Name: "Coordinator"},
		{Code: Accountant, Name: "Accountant"},
		{Code: Librarian, Name: "Librarian"},
		{Code: Receptionist, Name: "Receptionist"},
	}
}

func Known(code Code) bool {
	for _, item := range Catalog() {
		if item.Code == code {
			return true
		}
	}
	return false
}
