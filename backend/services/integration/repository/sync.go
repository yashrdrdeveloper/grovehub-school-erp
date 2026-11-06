package repository

import (
 "context"
 "grovehub/services/integration/entity"
)
type MailSyncRepository interface{
 Cursor(ctx context.Context,connectionID string)(string,error)
 SaveMessages(ctx context.Context,connectionID string,messages []entity.MailReference)error
 MarkDeleted(ctx context.Context,connectionID string,messageIDs []string)error
 SaveCursor(ctx context.Context,connectionID,cursor string)error
}
