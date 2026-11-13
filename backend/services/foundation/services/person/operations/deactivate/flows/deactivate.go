package flows

import (
    persondto "grovehub/services/foundation/dto/person"
    "grovehub/services/foundation/entity"
)

type Deactivate struct {
    Request persondto.DeactivatePersonRequest
    GID string
    Person entity.Person
    Contacts []entity.Contact
    UpdatedPerson entity.Person
}
