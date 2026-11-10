package repository

import (
	"context"
	"errors"

	"grovehub/services/foundation/entity"
)

var ErrContactAlreadyExists = errors.New("person contact already belongs to another person")

// @Repository entity=Contact
// ContactRepository exclusively owns persistence of foundation.person_contacts.
// Contact remains a Person-owned child; no ContactService or independent API exists.
type ContactRepository interface {
	FindByPersonGID(ctx context.Context, personGID string) ([]entity.Contact, error)
	FindByPersonGIDs(ctx context.Context, personGIDs []string) (map[string][]entity.Contact, error)
	FindPersonGIDsByContacts(ctx context.Context, contacts []entity.Contact) ([]string, error)
	ReplaceForPerson(ctx context.Context, personGID string, contacts []entity.Contact) error
}
