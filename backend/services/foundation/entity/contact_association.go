package entity

import "time"

// ForPerson attaches a contact to its owning Person while retaining the creation time.
func (contact Contact) ForPerson(gid GID, createdAt time.Time, updatedAt time.Time) Contact {
	contact.PersonGID = gid
	if contact.CreatedAt.IsZero() {
		contact.CreatedAt = createdAt
	}
	contact.UpdatedAt = updatedAt
	return contact
}
