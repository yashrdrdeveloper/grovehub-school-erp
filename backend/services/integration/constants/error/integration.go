package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
 Invalid = coreerrors.NewDefinition("INTEGRATION_INVALID","Invalid integration request")
 Forbidden = coreerrors.NewDefinition("INTEGRATION_FORBIDDEN","Integration access denied")
 NotFound = coreerrors.NewDefinition("INTEGRATION_NOT_FOUND","Integration connection not found")
 NotReady = coreerrors.NewDefinition("INTEGRATION_NOT_READY","Provider authorization not implemented")
 SyncBusy = coreerrors.NewDefinition("INTEGRATION_SYNC_BUSY","Mailbox sync is in progress")
 SyncLimit = coreerrors.NewDefinition("INTEGRATION_SYNC_LIMIT","Mailbox exceeds the safe per-run sync page limit")
 Failed = coreerrors.NewDefinition("INTEGRATION_FAILED","Integration operation failed")
)
