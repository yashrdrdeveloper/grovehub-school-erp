package helpers

import corenormalize "grovehub/libs/core/normalize"

// ValidID applies GrooveHub's standard UUIDv7 validation consistently.
func ValidID(raw string) bool {
    _,err:=corenormalize.UUIDv7(raw)
    return err==nil
}
