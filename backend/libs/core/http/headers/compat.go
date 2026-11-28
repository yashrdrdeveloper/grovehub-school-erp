// Package headers is a compatibility facade.
// Deprecated: use the root core/http package for new code.
package headers

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/internal/src/constants"
const (
 IdempotencyKey = constants.IdempotencyKey
 CorrelationID = constants.CorrelationID
 RequestID = constants.RequestID
 TraceID = constants.TraceID
)
