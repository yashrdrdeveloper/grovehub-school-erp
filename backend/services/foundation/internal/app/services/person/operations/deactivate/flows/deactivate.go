package flows

import (
    persondto "grovehub/services/foundation/client/person/dto"
    "grovehub/services/foundation/internal/app/entity"
)

type Deactivate struct {
    Request persondto.DeactivatePersonRequest
    GID string
    Person entity.Person
    Contacts []entity.Contact
    UpdatedPerson entity.Person
}
