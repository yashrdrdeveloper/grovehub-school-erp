package request

import (
	"bytes"
	"encoding/json"
)

// PatchString tracks presence independently of a nullable JSON string.
// Set=false means omitted; Set=true and Value=nil means explicit null.
type PatchString struct {
	Set   bool
	Value *string
}

func (field *PatchString) UnmarshalJSON(data []byte) error {
	field.Set = true
	field.Value = nil
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		return nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	field.Value = &value
	return nil
}

// PatchContacts tracks omitted, null, and supplied contacts independently.
type PatchContacts struct {
	Set   bool
	Value *Contacts
}

func (field *PatchContacts) UnmarshalJSON(data []byte) error {
	field.Set = true
	field.Value = nil
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		return nil
	}
	var value Contacts
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	field.Value = &value
	return nil
}

// MarshalJSON preserves omitted versus explicit-null fields on round trips.
func (request UpdatePerson) MarshalJSON() ([]byte, error) {
	fields := make(map[string]any, 4)
	if request.FirstName.Set { fields["firstName"] = request.FirstName.Value }
	if request.MiddleName.Set { fields["middleName"] = request.MiddleName.Value }
	if request.LastName.Set { fields["lastName"] = request.LastName.Value }
	if request.Contacts.Set { fields["contacts"] = request.Contacts.Value }
	return json.Marshal(fields)
}
