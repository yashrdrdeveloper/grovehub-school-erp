package helpers

import corenormalize "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/normalize"

// ValidID applies GrooveHub's standard UUIDv7 validation consistently.
func ValidID(raw string) bool {
    _,err:=corenormalize.UUIDv7(raw)
    return err==nil
}
