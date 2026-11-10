package repository

import (
	"context"

	"grovehub/services/foundation/entity"
)

/**
 * ContactRepository exclusively persists foundation.person_contacts rows.
 * Contacts are Person-owned children, not independently mutable resources.
 * Changes to contacts must be coordinated by the owning PersonService.
 */
// @Repository entity=Contact
type ContactRepository interface {
	FindByPersonGID(ctx context.Context, personGID string) ([]entity.Contact, error)
	FindByPersonGIDs(ctx context.Context, personGIDs []string) (map[string][]entity.Contact, error)
	FindPersonGIDsByContacts(ctx context.Context, contacts []entity.Contact) ([]string, error)
	ReplaceForPerson(ctx context.Context, personGID string, contacts []entity.Contact) error
}
