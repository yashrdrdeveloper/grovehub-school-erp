package schema

type Repository interface {
	CheckActivity() error
	CheckRejection() error
	CheckReceipt() error
	CheckOperator() error
	CheckSchoolReceipt() error
}
