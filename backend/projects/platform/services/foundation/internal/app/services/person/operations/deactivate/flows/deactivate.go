package flows

import (
    persondto "grovehub/services/foundation/internal/app/dto/person"
    "grovehub/services/foundation/internal/app/entity"
)

type Deactivate struct {
    Request persondto.DeactivatePersonRequest
    GID string
    Person entity.Person
    Contacts []entity.Contact
    UpdatedPerson entity.Person
}
