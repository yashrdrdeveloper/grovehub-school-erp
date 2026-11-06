package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
 Invalid = coreerrors.NewDefinition("INTEGRATION_INVALID","Invalid integration request")
 Forbidden = coreerrors.NewDefinition("INTEGRATION_FORBIDDEN","Integration access denied")
 NotFound = coreerrors.NewDefinition("INTEGRATION_NOT_FOUND","Integration connection not found")
 NotReady = coreerrors.NewDefinition("INTEGRATION_NOT_READY","Provider authorization not implemented")
 Failed = coreerrors.NewDefinition("INTEGRATION_FAILED","Integration operation failed")
)
