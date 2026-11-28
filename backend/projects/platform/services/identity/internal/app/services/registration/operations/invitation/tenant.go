package invitation

import (
 "errors"
 schoolgateway "grovehub/launcher/session"
)

// An unavailable registry cannot establish permanent ineligibility. Retain
// the job for retry. A verified stale or disabled proof must never be sent.
func tenantDeliveryDecision(state schoolgateway.TenantState, err error, issued int64) (string,error) {
 if err != nil { return "",err }
 if state.Generation < 1 { return "",errors.New("tenant generation unavailable") }
 if issued != state.Generation { return "INVALID_PROOF",nil }
 if !state.Active { return "ACCOUNT_INELIGIBLE",nil }
 return "",nil
}
