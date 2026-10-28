package contracts

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

type Direction string

const (
	DirectionAscending  Direction = "ASC"
	DirectionDescending Direction = "DESC"
)

const (
	DefaultRows = 25
	MaxRows     = 100
)

type Filter struct {
	Field    string   `json:"field"`
	Operator Operator `json:"operator"`
	Values   []string `json:"values,omitempty"`
}

type Sort struct {
	Field     string    `json:"field"`
	Direction Direction `json:"direction"`
}

type PageRequest struct {
	Rows   int    `json:"rows,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

type PageResponse struct {
	Rows       int    `json:"rows"`
	NextCursor string `json:"nextCursor,omitempty"`
	HasNext    bool   `json:"hasNext"`
}

type Request struct {
	SearchText string      `json:"searchText,omitempty"`
	Filters    []Filter    `json:"filters,omitempty"`
	Sort       []Sort      `json:"sort,omitempty"`
	PageInfo   PageRequest `json:"pageInfo"`
}

type Result[T any] struct {
	Items    []T          `json:"items"`
	PageInfo PageResponse `json:"pageInfo"`
}

type QueryFilter struct {
	Field    string
	Operator Operator
	Values   []string
}

type QuerySort struct {
	Field     string
	Direction Direction
}

type Cursor[TId any] struct {
	Value    string
	CursorId TId
}

type Query[TId any] struct {
	SearchTerms []string
	Filters     []QueryFilter
	Sort        QuerySort
	Cursor      *Cursor[TId]
	Limit       int
}

type QueryResult[T any] struct {
	Items   []T
	HasNext bool
}

type CursorPayload struct {
	Field     string    `json:"field"`
	Direction Direction `json:"direction"`
	Value     string    `json:"value"`
	CursorId  string    `json:"cursorId"`
}

type ScalarField struct {
	Expression string
	Cast       string
}

type FixedCondition struct {
	Expression string
	Value      any
}

type RelationField struct {
	Table           string
	Alias           string
	ParentCondition string
	ValueExpression string
	FixedConditions []FixedCondition
}

type SortField struct {
	Expression string
	Cast       string
}

type Config struct {
	SearchExpression string
	CursorExpression string
	ScalarFields     map[string]ScalarField
	RelationFields   map[string]RelationField
	SortFields       map[string]SortField
}
