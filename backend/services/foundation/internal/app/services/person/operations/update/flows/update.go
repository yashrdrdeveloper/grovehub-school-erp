package flows

import (
    persondto "grovehub/services/foundation/client/person/dto"
    "grovehub/services/foundation/internal/app/entity"
)

type Update struct {
    Request persondto.UpdatePersonRequest
    Validated persondto.UpdatePersonRequest
    GID string
    Person entity.Person
    Contacts []entity.Contact
    UpdatedPerson entity.Person
    UpdatedContacts []entity.Contact
    ContactsChanged bool
}
