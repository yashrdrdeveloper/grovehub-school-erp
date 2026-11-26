package validators

import (
 "errors"
 "strings"
 "unicode/utf8"

 "github.com/google/uuid"
)

func ValidateIdentity(tenantID, code, name string) error {
 if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(tenantID) != tenantID {
  return errors.New("invalid tenant scope")
 }
 if code == "" || strings.TrimSpace(code) != code || utf8.RuneCountInString(code) > 64 {
  return errors.New("invalid accounting code")
 }
 if name == "" || strings.TrimSpace(name) != name || utf8.RuneCountInString(name) > 255 {
  return errors.New("invalid accounting name")
 }
 return nil
}

func ValidateID(id string) error {
 parsed, err := uuid.Parse(id)
 if err != nil || parsed == uuid.Nil || parsed.String() != id {
  return errors.New("invalid accounting identifier")
 }
 return nil
}
