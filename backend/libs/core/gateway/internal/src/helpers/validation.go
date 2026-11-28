package helpers

import (
	"strings"

	"grovehub/libs/core/gateway/internal/src/constants"
)

// ValidKeyID checks the existing credential key identifier constraints.
func ValidKeyID(keyID string) bool {
	return keyID != "" && len(keyID) <= constants.MaxKeyIDLength && !strings.ContainsAny(keyID, ". \t\r\n,")
}

// ValidActorPair requires tenant and actor identifiers to be present together.
func ValidActorPair(tenantID, actorGID string) bool {
	return (tenantID == "") == (actorGID == "")
}
