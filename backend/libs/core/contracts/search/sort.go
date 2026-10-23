package search

type Direction string

const (
	DirectionAscending  Direction = "ASC"
	DirectionDescending Direction = "DESC"
)

type Sort struct {
	Field     string    `json:"field"`
	Direction Direction `json:"direction"`
}
