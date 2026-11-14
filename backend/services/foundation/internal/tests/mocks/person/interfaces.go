// Package personmocks declares Foundation-only mock targets.
// Mock implementations are generated into the Go overlay, never the app.
package personmocks

import (
    create "grovehub/services/foundation/internal/app/services/person/operations/create/strategies"
    getbygid "grovehub/services/foundation/internal/app/services/person/operations/getbygid/strategies"
    update "grovehub/services/foundation/internal/app/services/person/operations/update/strategies"
    search "grovehub/services/foundation/internal/app/services/person/operations/search/strategies"
    deactivate "grovehub/services/foundation/internal/app/services/person/operations/deactivate/strategies"
    protected "grovehub/services/foundation/internal/app/controller/person/protected"
    "grovehub/libs/core/database/transaction"
    "grovehub/libs/core/events/outbox"
    eventperson "grovehub/services/foundation/internal/app/events/producers/person"
    mapperperson "grovehub/services/foundation/internal/app/mappers/person"
    "grovehub/services/foundation/internal/app/repository/contact"
    personrepository "grovehub/services/foundation/internal/app/repository/person"
    personservice "grovehub/services/foundation/internal/app/services/person"
    personvalidators "grovehub/services/foundation/internal/app/validators/person"
    infrastructureclient "grovehub/services/infrastructure/client"
)

// @Mock
type PersonRepository interface { personrepository.PersonRepository }

// @Mock
type ContactRepository interface { contact.ContactRepository }

// @Mock
type PersonService interface { personservice.PersonService }

// @Mock
type PersonMapper interface { mapperperson.PersonMapper }

// @Mock
type PersonValidator interface { personvalidators.PersonValidator }

// @Mock
type ChangedProducer interface { eventperson.ChangedProducer }

// @Mock
type UnitOfWork interface { transaction.UnitOfWork }

// @Mock
type OutboxStore interface { outbox.Store }

// @Mock
type MessagingClient interface { infrastructureclient.MessagingClient }

// @Mock
type CreateStrategy interface { create.CreateStrategy }

// @Mock
type UpdateStrategy interface { update.UpdateStrategy }

// @Mock
type SearchStrategy interface { search.SearchStrategy }

// @Mock
type GetByGIDStrategy interface { getbygid.GetByGIDStrategy }

// @Mock
type DeactivateStrategy interface { deactivate.DeactivateStrategy }

// @Mock
type PersonController interface { protected.PersonController }
