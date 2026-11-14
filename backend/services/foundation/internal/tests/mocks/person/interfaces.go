// Package personmocks declares Foundation-only mock targets.
// Mock implementations are generated into the Go overlay, never the app.
package personmocks

import (
    "grovehub/libs/core/database/transaction"
    "grovehub/libs/core/events/outbox"
    "grovehub/services/foundation/internal/app/events/producers/person"
    "grovehub/services/foundation/internal/app/mappers/person"
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
type PersonMapper interface { person.PersonMapper }

// @Mock
type PersonValidator interface { personvalidators.PersonValidator }

// @Mock
type ChangedProducer interface { person.ChangedProducer }

// @Mock
type UnitOfWork interface { transaction.UnitOfWork }

// @Mock
type OutboxStore interface { outbox.Store }

// @Mock
type MessagingClient interface { infrastructureclient.MessagingClient }
