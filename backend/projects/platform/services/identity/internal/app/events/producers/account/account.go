package account
import ("context";"grovehub/services/identity/internal/app/entity")
type ChangedProducer interface{
 Enqueue(ctx context.Context,operation string,old,next *entity.Account)error
}
