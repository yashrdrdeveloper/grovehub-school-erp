package search

type Operator string

const (
	OperatorEquals      Operator = "EQUALS"
	OperatorNotEquals   Operator = "NOT_EQUALS"
	OperatorIn          Operator = "IN"
	OperatorNotIn       Operator = "NIN"
	OperatorExists      Operator = "EXISTS"
	OperatorNotExists   Operator = "NOT_EXISTS"
	OperatorGreaterThan Operator = "GT"
	OperatorGTE         Operator = "GTE"
	OperatorLessThan    Operator = "LT"
	OperatorLTE         Operator = "LTE"
	OperatorContains    Operator = "CONTAINS"
	OperatorStartsWith  Operator = "STARTS_WITH"
)
