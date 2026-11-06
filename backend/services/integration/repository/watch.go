package repository

import (
 "context"
 "time"
 "grovehub/services/integration/entity"
)
type WatchRepository interface{
 SaveWatch(ctx context.Context,connectionID,historyID string,expires time.Time)error
 GetGmailByMailbox(ctx context.Context,email string)([]entity.Connection,error)
 ExpiringWatches(ctx context.Context,before time.Time)([]entity.Connection,error)
}
