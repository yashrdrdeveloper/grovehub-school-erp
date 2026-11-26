package connection

import "errors"

// ErrMailDeliveryUncertain means dispatch may have reached the provider.
// The same delivery identifier must never be dispatched again automatically.
var ErrMailDeliveryUncertain = errors.New("mail delivery outcome uncertain")

// ErrMailDeliveryConflict indicates reuse of an identifier with different content.
var ErrMailDeliveryConflict = errors.New("mail delivery identifier conflict")
