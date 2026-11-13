package flows

import (
    persondto "grovehub/services/foundation/dto/person"
    "grovehub/services/foundation/entity"
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
