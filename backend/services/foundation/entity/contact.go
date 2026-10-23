package entity

type ContactType string

const (
	ContactTypeEmail ContactType = "EMAIL"
	ContactTypePhone ContactType = "PHONE"
)

type Contact struct {
	contactType     ContactType
	value           string
	normalizedValue string
	isPrimary       bool
	isVerified      bool
}

func NewContact(
	contactType ContactType,
	value string,
	normalizedValue string,
	isPrimary bool,
	isVerified bool,
) Contact {
	return Contact{
		contactType:     contactType,
		value:           value,
		normalizedValue: normalizedValue,
		isPrimary:       isPrimary,
		isVerified:      isVerified,
	}
}

func (contact Contact) GetType() ContactType {
	return contact.contactType
}

func (contact Contact) GetValue() string {
	return contact.value
}

func (contact Contact) GetNormalizedValue() string {
	return contact.normalizedValue
}

func (contact Contact) GetIsPrimary() bool {
	return contact.isPrimary
}

func (contact Contact) GetIsVerified() bool {
	return contact.isVerified
}
