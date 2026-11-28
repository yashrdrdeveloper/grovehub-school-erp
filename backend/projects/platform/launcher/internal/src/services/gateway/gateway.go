// Package gateway composes the external and internal gateway operations.
package gateway

import (
	"net/http"

	inprocess "grovehub/launcher/internal/src/services/gateway/operations/inprocess"
)

type Gateway interface {
 External() http.Handler
 Internal() inprocess.Operation
}

