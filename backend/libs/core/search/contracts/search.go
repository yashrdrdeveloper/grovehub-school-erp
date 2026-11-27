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
	Rows   int    `json:"rows,omitempty"` // legacy request alias
	PageNumber int `json:"pageNumber,omitempty"`
	RowSize int `json:"rowSize,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

type PageResponse struct {
	TotalRecords *int64 `json:"totalRecords,omitempty"`
	Rows       int    `json:"rows"`
	PageNumber int `json:"pageNumber,omitempty"`
	RowSize int `json:"rowSize,omitempty"`
	TotalPages *int64 `json:"totalPages,omitempty"`
	HasPrevious bool `json:"hasPrevious"`
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
	SecondaryValues []string
	CursorId TId
}

type Query[TId any] struct {
	SearchTerms []string
	Filters     []QueryFilter
	Sort        QuerySort
	SecondarySorts []QuerySort
	Cursor      *Cursor[TId]
	Limit       int
	PageNumber int
	Offset int
}

type QueryResult[T any] struct {
	TotalRecords *int64
	Items   []T
	HasNext bool
}

type CursorPayload struct {
	Field     string    `json:"field"`
	Direction Direction `json:"direction"`
	Value     string    `json:"value"`
	CursorId  string    `json:"cursorId"`
	SecondarySorts []Sort `json:"secondarySorts,omitempty"`
	SecondaryValues []string `json:"secondaryValues,omitempty"`
	PageNumber int `json:"pageNumber,omitempty"`
	RowSize int `json:"rowSize,omitempty"`
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
