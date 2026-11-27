package entity

type Persona string
const (
 L1 Persona = "L1"
 L2 Persona = "L2"
)

type Operator struct {
 ID string `json:"id"`
 Email string `json:"email"`
 Persona Persona `json:"persona"`
}

type Credential struct {
 Operator Operator
 Salt,Hash []byte
 Active bool
}
