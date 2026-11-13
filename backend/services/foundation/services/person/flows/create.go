// Package flows holds operation-local state shared by the Person strategies.
package flows

import (
    persondto "grovehub/services/foundation/dto/person"
    "grovehub/services/foundation/entity"
)

type Create struct {
    Request persondto.CreatePersonRequest
    Contacts []entity.Contact
    IdempotencyKey string
    RequestHash string
    Person entity.Person
}
