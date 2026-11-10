package impl

import (
    foundationclient "grovehub/services/foundation/client"
    foundationservices "grovehub/services/foundation/services"
)

// personClient forwards the public Foundation client to its owning service.
// Typed methods are generated from @Client/@Proxy on client.Client.
type personClient struct {
    personService foundationservices.PersonService
}

func NewPersonClient(personService foundationservices.PersonService) foundationclient.Client {
    return &personClient{personService: personService}
}
