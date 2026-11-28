package contract

// UseCase describes one supported context for a directive.
type UseCase struct {
	Name     string
	Syntax   string
	Category string
	Purpose  string
	Mode     string
	Status   Status
	Roadmap  string
	Target   Target
}
